package tests

import (
	"testing"

	"ngnix-superopt/pkg/asmopt"
	"ngnix-superopt/pkg/benchmarks"
	"ngnix-superopt/pkg/bolt"
	"ngnix-superopt/pkg/dpdk_fstack"
	"ngnix-superopt/pkg/hugepages"
	"ngnix-superopt/pkg/kernel_bypass"
	"ngnix-superopt/pkg/lifter"
	"ngnix-superopt/pkg/orchestrator"
	"ngnix-superopt/pkg/superopt"
	"ngnix-superopt/pkg/uop_tune"
	"ngnix-superopt/pkg/verifier"
)

func TestFullGoAssemblyPipeline(t *testing.T) {
	orch := orchestrator.NewOrchestrator("nginx")
	report, err := orch.RunFullPipeline()
	if err != nil {
		t.Fatalf("pipeline run failed: %v", err)
	}
	if report.EquivalenceProof.Equivalence != "SMT_VERIFIED_UNSAT" {
		t.Fatalf("expected SMT_VERIFIED_UNSAT, got %s", report.EquivalenceProof.Equivalence)
	}
	if len(report.AuditResults) == 0 {
		t.Fatalf("expected audit results")
	}
	for _, a := range report.AuditResults {
		if !a.Passed {
			t.Fatalf("audit check %s failed: %s", a.Name, a.Detail)
		}
	}
}

func TestSuperOptimizerAndVerifier(t *testing.T) {
	l := lifter.NewLifter()
	symbols := []string{"ngx_http_process_request_headers"}
	irs, err := l.LiftBinary("nginx", symbols)
	if err != nil {
		t.Fatalf("lifting failed: %v", err)
	}
	opt := superopt.NewSuperOptimizer()
	v := verifier.NewGatekeeperVerifier()
	for _, fn := range irs {
		optFn := opt.OptimizeFunction(fn)
		res := v.VerifyEquivalence(fn, optFn)
		if res.Equivalence != "SMT_VERIFIED_UNSAT" {
			t.Fatalf("verifier failed: %+v", res)
		}
	}
}

func TestAssemblyPrimitives(t *testing.T) {
	buf := []byte("X-Gateway-Trace-Id: abcdef123456")
	c1 := asmopt.AVX2ChecksumBytes(buf)
	if c1 == 0 {
		t.Fatalf("expected valid checksum")
	}

	pos := asmopt.FastStringMatchSearch(buf, []byte("Trace-Id"))
	if pos < 0 {
		t.Fatalf("expected to match substring in assembly")
	}

	zeroBuf := make([]byte, 64)
	for i := range zeroBuf {
		zeroBuf[i] = 0x55
	}
	asmopt.MemoryZeroAligned(zeroBuf)
	for i, b := range zeroBuf {
		if b != 0 {
			t.Fatalf("zeroBuf[%d] = %x, expected 0", i, b)
		}
	}
}

func TestSubsystems(t *testing.T) {
	boltOpt := bolt.NewBOLTOptimizer()
	layout := boltOpt.GenerateProfiledLayout([]string{"fn1", "fn2"})
	if len(layout.FunctionOrder) != 2 {
		t.Fatalf("bolt order failed")
	}

	hp := hugepages.NewHugePagesManager()
	if hp.Status() == "" {
		t.Fatalf("hugepages status empty")
	}

	tuner := uop_tune.NewMicroarchTuner()
	stats := tuner.AnalyzeLoopAlignment("test_fn")
	if stats.IPC <= 0 {
		t.Fatalf("invalid ipc")
	}

	fuzzer := verifier.NewFuzzer()
	ok, err := fuzzer.RunFuzzSuite(100)
	if !ok || err != nil {
		t.Fatalf("fuzz suite failed: %v", err)
	}

	bench := benchmarks.NewLiveBenchmark()
	base := bench.RunBaseline(100)
	super := bench.RunSuperopt(100)
	if base.TotalRequests != 100 || super.TotalRequests != 100 {
		t.Fatalf("benchmark request counts mismatch")
	}

	dpdk := dpdk_fstack.NewDPDKEngine()
	metrics := dpdk.GetMetrics()
	if metrics["mode"] != "userspace_tcp_pmd" {
		t.Fatalf("unexpected dpdk mode")
	}

	bypass := kernel_bypass.NewKernelBypassManager()
	info := bypass.InspectXDPProgram()
	if info["status"] != "ACTIVE" {
		t.Fatalf("xdp not active")
	}
}
