package main

import (
	"github.com/decomp/exp/bin"
	_ "github.com/decomp/exp/bin/elf" // register ELF decoder
	_ "github.com/decomp/exp/bin/pe"  // register PE decoder
	"github.com/lapsang-boys/pippi/cmd/pi-disasm/disasm"
	_ "github.com/lapsang-boys/pippi/cmd/pi-disasm/disasm/x86" // register 32- and 64-bit x86 disassemblers
	"github.com/mewkiz/pkg/jsonutil"
	"github.com/mewkiz/pkg/pathutil"
	"github.com/pkg/errors"
)

// A Database holds the state of a reverse engineering project.
type Database struct {
	// Addresses of instructions in the binary.
	InstAddrs []bin.Address
}

// parseDatabase parses the database of the given binary.
func parseDatabase(binPath string) (*Database, error) {
	dbPath := pathutil.TrimExt(binPath) + ".json"
	var instAddrs []bin.Address
	if err := jsonutil.ParseFile(dbPath, &instAddrs); err != nil {
		return nil, errors.WithStack(err)
	}
	db := &Database{
		InstAddrs: instAddrs,
	}
	return db, nil
}

// disasmBinary disassembles the given binary.
func disasmBinary(db *Database, arch bin.Arch, binPath string) ([]disasm.Instruction, error) {
	// Parse bianry file.
	file, err := bin.ParseFile(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Create disassembler based on machine architecture.
	dis, err := disasm.NewDisassembler(arch)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Disassemble instructions.
	var insts []disasm.Instruction
	for _, instAddr := range db.InstAddrs {
		data := file.Code(instAddr)
		inst, err := dis.DecodeInst(instAddr, data)
		if err != nil {
			//fmt.Fprintln(os.Stderr, hex.Dump(data))
			warn.Printf("unable to decode instruction at address %v", instAddr)
			continue
		}
		insts = append(insts, inst)
	}
	return insts, nil
}
