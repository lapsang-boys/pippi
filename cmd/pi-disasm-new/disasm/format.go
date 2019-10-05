// Note, the machine architecture format registration implementation of this
// package is heavily inspired by the image package of the Go standard library,
// which is governed by a BSD license.

package disasm

import (
	"github.com/decomp/exp/bin"
	"github.com/pkg/errors"
)

// RegisterFormat registers a machine architecture format for use by DecodeInst
// of a Disassembler. Name is the name of the format, like "x86_32" or
// "MIPS_32". Arch is is the machine architecture that identifies the format's
// Instruction Set Architecture (ISA).
func RegisterFormat(name string, arch bin.Arch, decodeInst func(addr bin.Address, buf []byte) (Instruction, error)) {
	formats = append(formats, format{name: name, arch: arch, decodeInst: decodeInst})
}

// formats is the list of registered formats.
var formats []format

// A format holds a machine architecture format's name, and details how to
// decode the instructions of its Instruction Set Architecture (ISA).
type format struct {
	// Name of the binary executable format.
	name string
	// Machine architecture that identifies the format's encoding.
	arch bin.Arch
	// decodeInst decodes the first instruction in buf.
	decodeInst func(addr bin.Address, buf []byte) (Instruction, error)
}

// DecodeInst decodes the first instruction in buf.
func (format *format) DecodeInst(addr bin.Address, buf []byte) (Instruction, error) {
	return format.decodeInst(addr, buf)
}

// NewDisassembler returns a new disassembler for the given machine
// architecture.
func NewDisassembler(arch bin.Arch) (Disassembler, error) {
	for _, format := range formats {
		if format.arch == arch {
			return &format, nil
		}
	}
	return nil, errors.Errorf("unknown machine architecture format %v;\n\ttip: remember to register a disassembler (e.g. import _ \".../disasm/x86\")", arch)
}
