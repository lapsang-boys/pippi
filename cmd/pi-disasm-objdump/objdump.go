package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/decomp/exp/bin"
	"github.com/pkg/errors"
)

// extractInstAddrs extracts the instruction addresses of the given binary using
// the objdump tool.
func extractInstAddrs(binPath string) ([]bin.Address, error) {
	cmd := exec.Command("objdump", "-d", binPath)
	buf, err := cmd.Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// $ objdump -d /usr/bin/ls
	//
	// /usr/bin/ls:     file format elf64-x86-64
	//
	// Disassembly of section .init:
	//
	// 0000000000004000 <.init>:
	//     4000:	f3 0f 1e fa          	endbr64
	//     4004:	48 83 ec 08          	sub    $0x8,%rsp
	//     4008:	48 8b 05 59 de 01 00 	mov    0x1de59(%rip),%rax
	var instAddrs []bin.Address
	s := bufio.NewScanner(bytes.NewReader(buf))
	for s.Scan() {
		line := s.Text()
		// Instruction lines are prefixed with space, filter other lines.
		//
		//     4000:	f3 0f 1e fa          	endbr64
		if !strings.HasPrefix(line, " ") {
			// skip line not starting with space.
			continue
		}
		pos := strings.IndexByte(line, ':')
		if pos == -1 {
			// skip line not containing colon.
			continue
		}
		line = "0x" + strings.TrimSpace(line[:pos])
		// Parse address.
		// TODO: check if we need to prefix with 0x before calling bin.Address.Set.
		var instAddr bin.Address
		if err := instAddr.Set(line); err != nil {
			return nil, errors.WithStack(err)
		}
		instAddrs = append(instAddrs, instAddr)
	}
	return instAddrs, nil
}
