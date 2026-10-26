"""
Integration test for DPDK and F-Stack Userspace TCP subsystem.
"""
import unittest
from dpdk_fstack.src.fstack_engine import FStackUserspaceTcpEngine, DpdkEalConfig

class TestDpdkFStackSubsystem(unittest.TestCase):
    def test_eal_initialization(self):
        engine = FStackUserspaceTcpEngine(DpdkEalConfig(core_mask="0x0e", hugepages_mb=4096))
        init_res = engine.init_eal()
        self.assertEqual(init_res["status"], "EAL_INITIALIZED")
        self.assertEqual(init_res["kernel_bypass"], "ACTIVE")

    def test_poll_burst_zero_syscalls(self):
        engine = FStackUserspaceTcpEngine()
        engine.init_eal()
        mbufs = engine.simulate_poll_burst(burst_count=128)
        self.assertEqual(len(mbufs), 128)
        self.assertEqual(engine.processed_packets, 128)

    def test_10x_performance_gain_profile(self):
        engine = FStackUserspaceTcpEngine()
        prof = engine.evaluate_10x_profile(target_cores=4)
        self.assertGreaterEqual(prof["speedup_vs_baseline_ratio"], 10.0)
        self.assertLess(prof["latency_p99_9_microseconds"], 500.0)
        self.assertEqual(prof["syscall_overhead_cycles"], 0)

if __name__ == "__main__":
    unittest.main()
