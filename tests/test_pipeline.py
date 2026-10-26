"""
Unit and Integration Tests for Pipeline Components using unittest.
Tests SMT gatekeeper on true equivalence and intentional divergence,
ensuring faulty candidates are rejected without exception.
"""
import unittest
import tempfile
from pipeline.verifier.ir import BasicBlock, MicroOp, OptimizationCandidate
from pipeline.verifier.gatekeeper import SmtGatekeeper
from pipeline.orchestrator import PipelineOrchestrator

class TestPipelineComponents(unittest.TestCase):
    def test_smt_gatekeeper_proven_equivalent(self):
        # Test valid equivalence: rax = rsi; rbx = 0 vs xor rbx, rbx
        orig_bb = BasicBlock(
            label="test_orig",
            instructions=[
                MicroOp("mov", "rax", ["rsi"]),
                MicroOp("mov", "rbx", ["0"]),
            ],
            reads_regs=["rsi"],
            writes_regs=["rax", "rbx"]
        )
        cand_bb = BasicBlock(
            label="test_cand",
            instructions=[
                MicroOp("mov", "rax", ["rsi"]),
                MicroOp("xor", "rbx", ["rbx"]),
            ],
            reads_regs=["rsi"],
            writes_regs=["rax", "rbx"]
        )
        cand = OptimizationCandidate("test_func", orig_bb, cand_bb, "xor zeroing", 1.2)
        gatekeeper = SmtGatekeeper()
        is_valid, msg, cex = gatekeeper.verify_equivalence(cand)
        self.assertTrue(is_valid)
        self.assertIn("PROVEN_EQUIVALENT", msg)
        self.assertIsNone(cex)

    def test_smt_gatekeeper_rejects_divergence(self):
        # Test faulty candidate: adds 1 instead of copy
        orig_bb = BasicBlock(
            label="test_orig_bug",
            instructions=[MicroOp("mov", "rax", ["rsi"])],
            reads_regs=["rsi"],
            writes_regs=["rax"]
        )
        cand_bb = BasicBlock(
            label="test_cand_bug",
            instructions=[
                MicroOp("mov", "rax", ["rsi"]),
                MicroOp("add", "rax", ["1"]),
            ],
            reads_regs=["rsi"],
            writes_regs=["rax"]
        )
        cand = OptimizationCandidate("test_bug", orig_bb, cand_bb, "buggy add", 1.0)
        gatekeeper = SmtGatekeeper()
        is_valid, msg, cex = gatekeeper.verify_equivalence(cand)
        self.assertFalse(is_valid)
        self.assertIn("REJECTED_DIVERGENCE", msg)
        self.assertIsNotNone(cex)

    def test_pipeline_orchestrator_execution(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            orchestrator = PipelineOrchestrator(tmp_dir)
            res_nginx = orchestrator.run_pipeline("nginx")
            self.assertEqual(res_nginx["target"], "nginx")
            self.assertGreater(res_nginx["stage_3_smt_gatekeeper"]["verified_equivalent_unsat"], 0)
            self.assertEqual(res_nginx["stage_5_testing"]["contour_1_fuzzing"]["verdict"], "PASSED_ALL_10M_MUTATIONS")
            self.assertEqual(res_nginx["stage_5_testing"]["contour_1_integration"]["passed"], 1284)

            res_envoy = orchestrator.run_pipeline("envoy")
            self.assertEqual(res_envoy["target"], "envoy")
            self.assertGreater(res_envoy["stage_3_smt_gatekeeper"]["verified_equivalent_unsat"], 0)
            self.assertEqual(res_envoy["stage_5_testing"]["contour_1_integration"]["passed"], 892)

if __name__ == "__main__":
    unittest.main()
