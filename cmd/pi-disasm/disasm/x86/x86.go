// Package x86 implements a disassembler for the 32- and 64-bit x86 machine
// architectures.
package x86

import (
	"github.com/decomp/exp/bin"
	"github.com/lapsang-boys/pippi/cmd/pi-disasm/disasm"
	"github.com/pkg/errors"
	"golang.org/x/arch/x86/x86asm"
)

func init() {
	// Register disassemblers for x86_32 and x86_64.
	disasm.RegisterFormat(bin.ArchX86_32.String(), bin.ArchX86_32, decodeInst32)
	disasm.RegisterFormat(bin.ArchX86_64.String(), bin.ArchX86_64, decodeInst64)
}

// decodeInst32 decodes the first instruction in buf.
func decodeInst32(addr bin.Address, buf []byte) (disasm.Instruction, error) {
	const mode = 32 // x86 (32-bit)
	inst, err := x86asm.Decode(buf, mode)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Instruction{Inst: inst, addr: addr}, nil
}

// decodeInst64 decodes the first instruction in buf.
func decodeInst64(addr bin.Address, buf []byte) (disasm.Instruction, error) {
	const mode = 64 // x86 (64-bit)
	inst, err := x86asm.Decode(buf, mode)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Instruction{Inst: inst, addr: addr}, nil
}

// Instruction is an x86 assembly instruction.
type Instruction struct {
	x86asm.Inst
	// Address of the instruction.
	addr bin.Address
}

// Addr returns the address of the instruction.
func (inst *Instruction) Addr() bin.Address {
	return inst.addr
}
