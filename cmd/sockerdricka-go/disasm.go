package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/decomp/exp/bin"
	_ "github.com/decomp/exp/bin/elf" // register ELF decoder
	_ "github.com/decomp/exp/bin/pe"  // register PE decoder
	"github.com/google/subcommands"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"golang.org/x/arch/x86/x86asm"
)

// disasmCmd is the command to disassemble a binary file from the command line.
type disasmCmd struct{}

func (*disasmCmd) Name() string {
	return "disasm"
}

func (*disasmCmd) Synopsis() string {
	return "disassemble binary file from command line"
}

func (*disasmCmd) Usage() string {
	const use = `
Disassemble binary file from command line.

Usage:
	disasm [OPTION]... FILE

Flags:
`
	return use[1:]
}

func (cmd *disasmCmd) SetFlags(f *flag.FlagSet) {
}

func (cmd *disasmCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	binPath := f.Arg(0)
	// Disassemble binary file and pretty-print to standard output.
	if err := disasm(binPath); err != nil {
		warn.Printf("disasm failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// disasm disassembles the given binary file, pretty-printing its assembly
// instructions to standard output.
func disasm(binPath string) error {
	// Parse binary file.
	file, err := bin.ParseFile(binPath)
	if err != nil {
		return errors.WithStack(err)
	}
	// Processor mode (16-, 32-, or 64-bit exection mode).
	var mode int
	switch file.Arch {
	case bin.ArchX86_32:
		mode = 32
	case bin.ArchX86_64:
		mode = 64
	}
	for _, sect := range file.Sections {
		fmt.Println("sect:", sect.Name)
		if sect.Perm&bin.PermX != 0 {
			valid, err := shingledDisasm(mode, sect.Data, uint64(sect.Addr))
			if err != nil {
				return errors.WithStack(err)
			}
			// Pretty-print to standard output.
			pretty.Println(valid)
		}
	}
	return nil
}

// shingledDisasm determines the superset of valid instructions in the given
// section data using the Shingled Graph Disassembly method. The mode argument
// specifies the processor mode (16-, 32-, or 64-bit exection mode).
func shingledDisasm(mode int, sectData []byte, sectAddr uint64) (valid []bool, err error) {
	n := len(sectData)
	valid = make([]bool, n)
	for i := range sectData {
		valid[i] = isValidEnc(mode, sectData[i:])
	}
	visited := make([]bool, n)
	q := &queue{
		buf: make([]uint64, 0, n),
	}
	for i := range sectData {
		if visited[i] || !valid[i] {
			continue
		}
		q.reset()
		visited[i] = true
		q.push(uint64(i))
		prune := false
		for !q.empty() {
			off := q.pop()
			if !valid[off] {
				prune = true
				continue
			}
			instAddr := sectAddr + off
			inst, err := x86asm.Decode(sectData[off:], mode)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			fmt.Printf("%08X\t%v\n", instAddr, inst)
			if hasFallthrough(inst) {
				nextOff := off + uint64(inst.Len)
				if nextOff < 0 || nextOff >= uint64(len(sectData)) {
					continue
				}
				if !visited[nextOff] {
					visited[nextOff] = true
					q.push(nextOff)
				}
			}
			for _, branchAddr := range branches(inst, instAddr) {
				branchOff := branchAddr - sectAddr
				if branchOff < 0 || branchOff >= uint64(len(sectData)) {
					continue
				}
				if !visited[branchOff] {
					visited[branchOff] = true
					q.push(branchOff)
				}
			}
		}
		if prune {
			for _, off := range q.buf {
				valid[off] = false
			}
		}
	}
	return valid, nil
}

// branches returns the addresses of the branches of the given instruction.
func branches(inst x86asm.Inst, instAddr uint64) []uint64 {
	switch inst.Op {
	// Loop terminators.
	case x86asm.LOOP, x86asm.LOOPE, x86asm.LOOPNE:
	// Conditional jump terminators.
	case x86asm.JA, x86asm.JAE, x86asm.JB, x86asm.JBE, x86asm.JCXZ, x86asm.JE, x86asm.JECXZ, x86asm.JG, x86asm.JGE, x86asm.JL, x86asm.JLE, x86asm.JNE, x86asm.JNO, x86asm.JNP, x86asm.JNS, x86asm.JO, x86asm.JP, x86asm.JRCXZ, x86asm.JS:
	// Unconditional jump terminators.
	case x86asm.JMP:
	// Return terminators.
	case x86asm.RET:
		return nil
	// Call instruction (not terminator, but includes branch other than
	// fallthrough).
	case x86asm.CALL:
	default:
		return nil
	}
	nextAddr := instAddr + uint64(inst.Len)
	switch arg := inst.Args[0].(type) {
	//case x86asm.Imm:
	case x86asm.Rel:
		return []uint64{nextAddr + uint64(arg)}
	default:
		// TODO: get branch target by symbolic exectuion.
		return nil
	}
}

// hasFallthrough reports whether the given instruction has a fall-through
// control flow semantic.
func hasFallthrough(inst x86asm.Inst) bool {
	switch inst.Op {
	case x86asm.JMP, x86asm.RET:
		return false
	}
	return true
}

// isValidEnc reports whether the start of the given data encodes a valid
// instruction.
func isValidEnc(mode int, data []byte) bool {
	_, err := x86asm.Decode(data, mode)
	return err == nil
}

type queue struct {
	buf []uint64
	i   int
}

func (q *queue) empty() bool {
	return q.i == len(q.buf)
}

func (q *queue) push(elem uint64) {
	q.buf = append(q.buf, elem)
}

func (q *queue) pop() uint64 {
	elem := q.buf[q.i]
	q.i++
	return elem
}

func (q *queue) reset() {
	q.buf = q.buf[:0]
	q.i = 0
}
