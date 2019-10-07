package disasm

import (
	"fmt"

	"github.com/decomp/exp/bin"
)

// Disassembler is a disassembler for a given machine architecture.
type Disassembler interface {
	// DecodeInst decodes the first instruction in buf.
	DecodeInst(addr bin.Address, buf []byte) (Instruction, error)
}

// Instruction is an assembly instruction of a given machine architecture.
type Instruction interface {
	fmt.Stringer
	// Addr returns the address of the instruction.
	Addr() bin.Address
}
