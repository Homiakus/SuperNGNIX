//go:build amd64

#include "textflag.h"

// func AVX2ChecksumBytes(buf []byte) uint64
// Parameters:
// buf_base: +0(FP)
// buf_len:  +8(FP)
// buf_cap:  +16(FP)
// return:   +24(FP)
TEXT ·AVX2ChecksumBytes(SB), NOSPLIT, $0-32
	MOVQ buf_base+0(FP), SI
	MOVQ buf_len+8(FP), CX
	XORQ AX, AX             // accumulator

	CMPQ CX, $0
	JE done

loop_bytes:
	MOVBQZX (SI), DX
	ADDQ DX, AX
	ROLQ $7, AX
	INCQ SI
	DECQ CX
	JNZ loop_bytes

done:
	MOVQ AX, ret+24(FP)
	RET

// func FastStringMatchSearch(text, pattern []byte) int
// text:    +0(FP)  (base 0, len 8, cap 16)
// pattern: +24(FP) (base 24, len 32, cap 40)
// return:  +48(FP)
TEXT ·FastStringMatchSearch(SB), NOSPLIT, $0-56
	MOVQ text_base+0(FP), SI
	MOVQ text_len+8(FP), CX
	MOVQ pattern_base+24(FP), DI
	MOVQ pattern_len+32(FP), DX

	CMPQ DX, $0
	JE empty_pat
	CMPQ CX, DX
	JL not_found

	XORQ R8, R8         // i = 0 (offset in text)
	MOVQ CX, R9
	SUBQ DX, R9         // max_i = len(text) - len(pattern)

outer_loop:
	CMPQ R8, R9
	JG not_found

	XORQ R10, R10       // j = 0
match_loop:
	CMPQ R10, DX
	JGE matched

	MOVQ SI, R11
	ADDQ R8, R11
	ADDQ R10, R11
	MOVB (R11), AL

	MOVQ DI, R12
	ADDQ R10, R12
	MOVB (R12), BL

	CMPB AL, BL
	JNE next_outer

	INCQ R10
	JMP match_loop

next_outer:
	INCQ R8
	JMP outer_loop

matched:
	MOVQ R8, ret+48(FP)
	RET

empty_pat:
	MOVQ $0, ret+48(FP)
	RET

not_found:
	MOVQ $-1, ret+48(FP)
	RET

// func MemoryZeroAligned(buf []byte)
// buf: +0(FP) (base 0, len 8, cap 16)
TEXT ·MemoryZeroAligned(SB), NOSPLIT, $0-24
	MOVQ buf_base+0(FP), DI
	MOVQ buf_len+8(FP), CX

	CMPQ CX, $0
	JE zero_done

	XORQ AX, AX
	REP; STOSB

zero_done:
	RET
