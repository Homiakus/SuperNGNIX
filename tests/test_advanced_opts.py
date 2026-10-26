"""
Enhanced Test Suite covering all 6 advanced optimization subsystems.
"""
import unittest
from pipeline.kernel_bypass.io_uring_xdp import IoUringEngine, EbpfXdpFilter
from pipeline.hugepages.hugepages import HugePagesOptimizer
from pipeline.bolt.bolt_optimizer import BoltLayoutOptimizer
from pipeline.resctrl.cache_cat import CacheAllocationManager
from pipeline.uop_tune.microarch_tuner import MicroarchitecturePipelineTuner
from pipeline.verifier.ir import BasicBlock, MicroOp

class TestAdvancedOptimizations(unittest.TestCase):
    def test_io_uring_syscall_elimination(self):
        engine = IoUringEngine()
        res = engine.benchmark_syscall_elimination(target_rps=200000)
        self.assertEqual(res["syscall_reduction_pct"], 100.0)
        self.assertGreater(res["throughput_boost_pct"], 15.0)

    def test_hugepages_tlb_reduction(self):
        opt = HugePagesOptimizer()
        tlb = opt.evaluate_tlb_efficiency(active_connections=10000)
        self.assertGreater(tlb["tlb_entry_reduction_factor"], 100)
        self.assertLess(tlb["optimized_tlb_miss_rate_pct"], 1.0)

    def test_bolt_layout_optimization(self):
        bolt = BoltLayoutOptimizer()
        cmd = bolt.build_bolt_command("nginx.bin", "perf.data", "nginx.bolt")
        self.assertIn("llvm-bolt", cmd)
        self.assertIn("-reorder-blocks=ext-tsp", cmd)
        stats = bolt.estimate_bolt_i_cache_improvement()
        self.assertLess(stats["i_cache_miss_reduction_pct"], -20.0)

    def test_cache_allocation_resctrl(self):
        cat = CacheAllocationManager()
        sh = cat.generate_resctrl_setup_commands()
        self.assertIn("resctrl", sh)
        eval_res = cat.evaluate_cache_isolation_benefit()
        self.assertEqual(eval_res["noisy_neighbor_l3_evictions"], "ELIMINATED")

    def test_microarch_pipeline_tuning(self):
        tuner = MicroarchitecturePipelineTuner()
        bb = BasicBlock("test", [
            MicroOp("cmp", "rax", ["0"]),
            MicroOp("je", "", ["target"]),
            MicroOp("mov", "rbx", ["1"])
        ])
        tune = tuner.align_basic_block_for_uop_cache(bb)
        self.assertEqual(tune["macro_fused_uops"], 1)
        self.assertTrue(tune["fits_in_single_dsb_line"])

if __name__ == "__main__":
    unittest.main()
