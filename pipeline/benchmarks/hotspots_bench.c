#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <immintrin.h>
#include <stdint.h>

// -------------------------------------------------------------
// 1. HTTP HEADER PARSER BENCHMARK
// Target: ngx_http_parse_header_line
// -------------------------------------------------------------

// Baseline: Scalar loop with branch cascade (GCC/Clang -O3 equivalent)
__attribute__((noinline))
int parse_header_baseline(const char *buf, size_t len) {
    int matches = 0;
    for (size_t i = 0; i < len; ++i) {
        char c = buf[i];
        if (c == '\n') {
            matches++;
        } else if (c == '\r') {
            matches++;
        }
    }
    return matches;
}

// LLM Super-Opt: AVX2 32-byte SIMD vector scan + bitmask + tzcnt
__attribute__((noinline))
int parse_header_avx2(const char *buf, size_t len) {
    int matches = 0;
    size_t i = 0;
    __m256i v_lf = _mm256_set1_epi8('\n');
    __m256i v_cr = _mm256_set1_epi8('\r');

    for (; i + 32 <= len; i += 32) {
        __m256i chunk = _mm256_loadu_si256((const __m256i*)(buf + i));
        __m256i cmp_lf = _mm256_cmpeq_epi8(chunk, v_lf);
        __m256i cmp_cr = _mm256_cmpeq_epi8(chunk, v_cr);
        __m256i combined = _mm256_or_si256(cmp_lf, cmp_cr);
        uint32_t mask = (uint32_t)_mm256_movemask_epi8(combined);
        matches += _mm_popcnt_u32(mask);
    }
    for (; i < len; ++i) {
        char c = buf[i];
        if (c == '\n' || c == '\r') matches++;
    }
    return matches;
}

// -------------------------------------------------------------
// 2. ENVOY FILTER DISPATCH BENCHMARK
// Target: envoy_filter_vtable_dispatch vs Devirtualized Direct Call
// -------------------------------------------------------------

struct Filter {
    int (*execute)(struct Filter *self, int ctx);
    int tag;
};

int filter_work(struct Filter *self, int ctx) {
    return ctx * 3 + self->tag;
}

// Baseline: Indirect dynamic call through function pointer / vtable
__attribute__((noinline))
int dispatch_baseline(struct Filter *f, int iterations) {
    int acc = 0;
    for (int i = 0; i < iterations; ++i) {
        acc += f->execute(f, i);
    }
    return acc;
}

// LLM Super-Opt: Speculative devirtualization & direct call inlining with guard
__attribute__((noinline))
int dispatch_devirtualized(struct Filter *f, int iterations) {
    int acc = 0;
    for (int i = 0; i < iterations; ++i) {
        if (__builtin_expect(f->execute == filter_work, 1)) {
            acc += (i * 3 + f->tag); // Direct inlined execution
        } else {
            acc += f->execute(f, i);
        }
    }
    return acc;
}

// -------------------------------------------------------------
// 3. ROUTING HASH TABLE BENCHMARK
// Target: ngx_hash_find (Instruction scheduling & prefetch)
// -------------------------------------------------------------

#define HASH_SIZE 4096
struct HashNode {
    uint32_t key;
    uint32_t val;
    struct HashNode *next;
};

struct HashNode nodes[HASH_SIZE];

void init_hash() {
    for (int i = 0; i < HASH_SIZE; ++i) {
        nodes[i].key = i * 17;
        nodes[i].val = i ^ 0x5a5a;
        nodes[i].next = (i < HASH_SIZE - 1) ? &nodes[i + 1] : NULL;
    }
}

__attribute__((noinline))
uint32_t hash_find_baseline(uint32_t *keys, int count) {
    uint32_t sum = 0;
    for (int i = 0; i < count; ++i) {
        uint32_t idx = (keys[i] ^ (keys[i] >> 4)) % HASH_SIZE;
        struct HashNode *curr = &nodes[idx];
        while (curr) {
            if (curr->key == keys[i]) {
                sum += curr->val;
                break;
            }
            curr = curr->next;
        }
    }
    return sum;
}

__attribute__((noinline))
uint32_t hash_find_optimized(uint32_t *keys, int count) {
    uint32_t sum = 0;
    for (int i = 0; i < count; ++i) {
        // Prefetch next key slot to hide L1/L2 latency (Memory-Level Parallelism)
        if (i + 4 < count) {
            uint32_t next_idx = (keys[i + 4] ^ (keys[i + 4] >> 4)) % HASH_SIZE;
            __builtin_prefetch(&nodes[next_idx], 0, 1);
        }
        uint32_t idx = (keys[i] ^ (keys[i] >> 4)) % HASH_SIZE;
        struct HashNode *curr = &nodes[idx];
        while (curr) {
            if (curr->key == keys[i]) {
                sum += curr->val;
                break;
            }
            curr = curr->next;
        }
    }
    return sum;
}

static inline double get_time_sec() {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return ts.tv_sec + ts.tv_nsec * 1e-9;
}

int main() {
    printf("{\"benchmark_suite\": \"Hardware Micro-Benchmark (Core i5-9300H AVX2)\",\n");

    // 1. HTTP Parser Test (Scenario A)
    size_t header_len = 1024 * 1024 * 64; // 64 MB buffer of headers
    char *hdr_buf = (char*)malloc(header_len);
    for (size_t i = 0; i < header_len; ++i) {
        if (i % 80 == 78) hdr_buf[i] = '\r';
        else if (i % 80 == 79) hdr_buf[i] = '\n';
        else hdr_buf[i] = 'A' + (i % 26);
    }

    double t0 = get_time_sec();
    volatile int r1 = parse_header_baseline(hdr_buf, header_len);
    double t_parser_base = get_time_sec() - t0;

    t0 = get_time_sec();
    volatile int r2 = parse_header_avx2(hdr_buf, header_len);
    double t_parser_opt = get_time_sec() - t0;

    // 2. Envoy Filter Dispatch (Scenario B)
    struct Filter f;
    f.execute = filter_work;
    f.tag = 42;
    int dispatch_iters = 50000000;

    t0 = get_time_sec();
    volatile int d1 = dispatch_baseline(&f, dispatch_iters);
    double t_disp_base = get_time_sec() - t0;

    t0 = get_time_sec();
    volatile int d2 = dispatch_devirtualized(&f, dispatch_iters);
    double t_disp_opt = get_time_sec() - t0;

    // 3. Routing Hash Table (Scenario C)
    init_hash();
    int key_count = 5000000;
    uint32_t *keys = (uint32_t*)malloc(sizeof(uint32_t) * key_count);
    for (int i = 0; i < key_count; ++i) {
        keys[i] = (rand() % HASH_SIZE) * 17;
    }

    t0 = get_time_sec();
    volatile uint32_t h1 = hash_find_baseline(keys, key_count);
    double t_hash_base = get_time_sec() - t0;

    t0 = get_time_sec();
    volatile uint32_t h2 = hash_find_optimized(keys, key_count);
    double t_hash_opt = get_time_sec() - t0;

    printf(" \"parser\": {\"baseline_sec\": %f, \"optimized_sec\": %f, \"speedup\": %f},\n",
           t_parser_base, t_parser_opt, t_parser_base / t_parser_opt);
    printf(" \"envoy_dispatch\": {\"baseline_sec\": %f, \"optimized_sec\": %f, \"speedup\": %f},\n",
           t_disp_base, t_disp_opt, t_disp_base / t_disp_opt);
    printf(" \"routing_hash\": {\"baseline_sec\": %f, \"optimized_sec\": %f, \"speedup\": %f}\n}\n",
           t_hash_base, t_hash_opt, t_hash_base / t_hash_opt);

    free(hdr_buf);
    free(keys);
    return 0;
}
