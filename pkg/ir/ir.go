package ir

import "fmt"

type OpType string

const (
	OpMov    OpType = "MOV"
	OpAdd    OpType = "ADD"
	OpSub    OpType = "SUB"
	OpMul    OpType = "MUL"
	OpLea    OpType = "LEA"
	OpJmp    OpType = "JMP"
	OpCmp    OpType = "CMP"
	OpNop    OpType = "NOP"
	OpVmov   OpType = "VMOVAPS"
	OpVxor   OpType = "VXORPS"
	OpPrefetch OpType = "PREFETCH"
)

type Instruction struct {
	Address uint64
	Op      OpType
	Args    []string
	RawAsm  string
}

func (i Instruction) String() string {
	if i.RawAsm != "" {
		return fmt.Sprintf("0x%x: %s", i.Address, i.RawAsm)
	}
	return fmt.Sprintf("0x%x: %s %v", i.Address, i.Op, i.Args)
}

type BasicBlock struct {
	ID           string
	Instructions []Instruction
	Predecessors []string
	Successors   []string
	ExecutionCount int64
}

type FunctionIR struct {
	Name        string
	Hotspot     bool
	Blocks      map[string]*BasicBlock
	BlockOrder  []string
	EntryBlock  string
}

func NewFunctionIR(name string) *FunctionIR {
	return &FunctionIR{
		Name:       name,
		Blocks:     make(map[string]*BasicBlock),
		BlockOrder: make([]string, 0),
	}
}

func (f *FunctionIR) AddBlock(b *BasicBlock) {
	f.Blocks[b.ID] = b
	f.BlockOrder = append(f.BlockOrder, b.ID)
	if f.EntryBlock == "" {
		f.EntryBlock = b.ID
	}
}
