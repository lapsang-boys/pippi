package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/subcommands"
	"github.com/lapsang-boys/pippi/cmd/alvriket/binpbx"
	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	// Default bin gRPC address to listen on.
	binGRPCAddr = ":1200"
	// Default disasm gRPC address to listen on.
	disasmGRPCAddr = ":1300"
	// Extension of binary files.
	ext = ".bin"
)

// serverCmd is the command to launch a gRPC server processing disassemble
// binary requests.
type serverCmd struct {
	// bin gRPC address to listen on.
	BinAddr string
	// disasm gRPC address to listen on.
	DisasmAddr string
}

func (*serverCmd) Name() string {
	return "server"
}

func (*serverCmd) Synopsis() string {
	return "launch gRPC server"
}

func (*serverCmd) Usage() string {
	const use = `
Reply to disassemble binary file request.

Usage:
	server [OPTION]...

Flags:
`
	return use[1:]
}

func (cmd *serverCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.BinAddr, "addr_bin", binGRPCAddr, "bin gRPC address to listen on")
	f.StringVar(&cmd.DisasmAddr, "addr_disasm", disasmGRPCAddr, "disasm gRPC address to connect to")
}

func (cmd *serverCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := listen(cmd.BinAddr, cmd.DisasmAddr); err != nil {
		warn.Printf("listen failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// listen listens on the given gRPC address for incoming requests to disassemble
// binary files.
func listen(binAddr, disasmAddr string) error {
	dbg.Printf("listening on %q", disasmAddr)
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return errors.WithStack(err)
	}
	cacheDir = filepath.Join(cacheDir, "pippi")
	// Launch gRPC server.
	l, err := net.Listen("tcp", disasmAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	server := grpc.NewServer()
	// Register disassembler service.
	disasmpb.RegisterDisassemblerServer(server, &disasmServer{cacheDir: cacheDir, binAddr: binAddr})
	if err := server.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// disasmServer implements disasmpb.DisassembleServer.
type disasmServer struct {
	// Cache directory of pippi.
	cacheDir string
	// bin gRPC address.
	binAddr string
}

// Disassemble disassembles the given binary file.
func (s *disasmServer) Disassemble(ctx context.Context, req *disasmpb.DisassembleRequest) (*disasmpb.DisassembleReply, error) {
	if err := validateID(req.BinId); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("disassembling ID %q", req.BinId)
	// Read file contents.
	binName := req.BinId + ext
	binPath := filepath.Join(s.cacheDir, req.BinId, binName)
	binData, err := ioutil.ReadFile(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Parse binary file.
	file, err := binpbx.ParseFile(s.binAddr, req.BinId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Send reply.
	reply := &disasmpb.DisassembleReply{}
	// Processor mode (16-, 32-, or 64-bit exection mode).
	var mode int
	switch file.Arch {
	case binpb.Arch_X86_32:
		mode = 32
	case binpb.Arch_X86_64:
		mode = 64
	}
	for _, sect := range file.Sections {
		if !permContains(sect.Perms, binpb.Perm_X) {
			continue
		}
		sectData := binData[sect.Offset : sect.Offset+sect.Length]
		valid, err := shingledDisasm(mode, sectData, sect.Addr)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		var validOffs []uint64
		for off := range valid {
			if valid[off] {
				validOffs = append(validOffs, uint64(off))
			}
		}
		execSect := &disasmpb.DisassembleSection{
			Section:      sect,
			ValidOffsets: validOffs,
		}
		reply.ExecSections = append(reply.ExecSections, execSect)
	}
	// Sections.
	return reply, nil
}

// permContains reports whether the given slice of permissions contains the
// specified permission.
func permContains(perms []binpb.Perm, perm binpb.Perm) bool {
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

// validateID validates the given binary ID.
func validateID(id string) error {
	if sha256.Size*2 != len(id) {
		return errors.Errorf("invalid length of binary ID; expected %d, got %d", sha256.Size*2, len(id))
	}
	s := strings.ToLower(id)
	if s != id {
		return errors.Errorf("invalid binary ID; expected lowercase, got %q", id)
	}
	const hex = "0123456789abcdef"
	for _, r := range id {
		if !strings.ContainsRune(hex, r) {
			return errors.Errorf("invalid rune in binary ID; expected hexadecimal digit, got %q", r)
		}
	}
	return nil
}
