// Package x86 implements a disassembler for the 32- and 64-bit x86 machine
// architectures.
package x86

import (
	"github.com/decomp/exp/bin"
	"github.com/lapsang-boys/pippi/cmd/pi-disasm-new/disasm"
	"github.com/pkg/errors"
	"golang.org/x/arch/x86/x86asm"
)

func init() {
	// Register disassemblers for x86_32 and x86_64.
	disasm.RegisterFormat(bin.ArchX86_32.String(), bin.ArchX86_32, decodeInst32)
	disasm.RegisterFormat(bin.ArchX86_64.String(), bin.ArchX86_64, decodeInst64)
}

// decodeInst32 decodes the first instruction in buf.
func decodeInst32(buf []byte) (disasm.Instruction, error) {
	const mode = 32 // x86 (32-bit)
	inst, err := x86asm.Decode(buf, mode)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return inst, nil
}

// decodeInst64 decodes the first instruction in buf.
func decodeInst64(buf []byte) (disasm.Instruction, error) {
	const mode = 64 // x86 (64-bit)
	inst, err := x86asm.Decode(buf, mode)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return inst, nil
}
