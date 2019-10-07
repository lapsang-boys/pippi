// Package mips implements a disassembler for the 32-bit MIPS machine
// architecture.
package mips

import (
	"github.com/decomp/exp/bin"
	"github.com/lapsang-boys/pippi/cmd/pi-disasm/disasm"
	"github.com/mewmew/mips"
	"github.com/pkg/errors"
)

func init() {
	// Register disassemblers for MIPS_32.
	disasm.RegisterFormat(bin.ArchMIPS_32.String(), bin.ArchMIPS_32, decodeInst32)
}

// decodeInst32 decodes the first instruction in buf.
func decodeInst32(addr bin.Address, buf []byte) (disasm.Instruction, error) {
	inst, err := mips.Decode(buf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Instruction{Inst: inst, addr: addr}, nil
}

// Instruction is a MIPS assembly instruction.
type Instruction struct {
	mips.Inst
	// Address of the instruction.
	addr bin.Address
}

// Addr returns the address of the instruction.
func (inst *Instruction) Addr() bin.Address {
	return inst.addr
}
