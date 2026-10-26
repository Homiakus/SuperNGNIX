"""
Instruction IR and Semantic Model for Assembly Blocks.
Provides intermediate representations for x86-64 instructions to model register/memory state in Z3.
"""
from dataclasses import dataclass, field
from typing import List, Dict, Any, Optional

@dataclass
class MicroOp:
    mnemonic: str
    dst: str
    src: List[str] = field(default_factory=list)
    raw_asm: str = ""

@dataclass
class BasicBlock:
    label: str
    instructions: List[MicroOp] = field(default_factory=list)
    reads_regs: List[str] = field(default_factory=list)
    writes_regs: List[str] = field(default_factory=list)
    reads_mem: bool = False
    writes_mem: bool = False

@dataclass
class OptimizationCandidate:
    function_name: str
    original_block: BasicBlock
    candidate_block: BasicBlock
    rationale: str
    speedup_est: float
