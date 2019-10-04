package disasm

import (
	"fmt"
)

// Disassembler is a disassembler for a given machine architecture.
type Disassembler interface {
	// DecodeInst decodes the first instruction in buf.
	DecodeInst(buf []byte) (Instruction, error)
}

// Instruction is an assembly instruction of a given machine architecture.
type Instruction interface {
	fmt.Stringer
}
