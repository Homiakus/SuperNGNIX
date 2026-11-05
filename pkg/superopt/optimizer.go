package superopt

import (
	"strings"

	"ngnix-superopt/pkg/ir"
)

type SuperOptimizer struct{}

func NewSuperOptimizer() *SuperOptimizer {
	return &SuperOptimizer{}
}

func (s *SuperOptimizer) OptimizeFunction(fn *ir.FunctionIR) *ir.FunctionIR {
	optFn := ir.NewFunctionIR(fn.Name)
	optFn.Hotspot = fn.Hotspot

	for _, bbID := range fn.BlockOrder {
		bb := fn.Blocks[bbID]
		optBB := &ir.BasicBlock{
			ID:             bb.ID,
			Predecessors:   append([]string(nil), bb.Predecessors...),
			Successors:     append([]string(nil), bb.Successors...),
			ExecutionCount: bb.ExecutionCount,
		}

		// Apply peephole and microarchitectural transforms
		for _, instr := range bb.Instructions {
			raw := strings.TrimSpace(instr.RawAsm)

			// Transform 1: zeroing idiom (mov $0, %eax -> xor %eax, %eax)
			if strings.HasPrefix(raw, "mov $0x0,") || strings.HasPrefix(raw, "mov $0,") {
				parts := strings.Split(raw, ",")
				if len(parts) == 2 {
					reg := strings.TrimSpace(parts[1])
					optBB.Instructions = append(optBB.Instructions, ir.Instruction{
						Address: instr.Address,
						Op:      ir.OpSub,
						RawAsm:  "xor " + reg + ", " + reg,
					})
					continue
				}
			}

			// Transform 2: add $1, %reg -> inc %reg (or lean lea)
			if strings.HasPrefix(raw, "add $0x1,") || strings.HasPrefix(raw, "add $1,") {
				parts := strings.Split(raw, ",")
				if len(parts) == 2 {
					reg := strings.TrimSpace(parts[1])
					optBB.Instructions = append(optBB.Instructions, ir.Instruction{
						Address: instr.Address,
						Op:      ir.OpAdd,
						RawAsm:  "inc " + reg,
					})
					continue
				}
			}

			// Pass-through
			optBB.Instructions = append(optBB.Instructions, instr)
		}

		optFn.AddBlock(optBB)
	}

	return optFn
}
