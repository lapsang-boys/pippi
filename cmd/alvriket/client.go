package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/lapsang-boys/pippi/alvriket"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// clientCmd is the command to launch a gRPC server processing parse ELF
// requests.
type clientCmd struct {
	// gRPC address to connect to.
	Addr string
}

func (*clientCmd) Name() string {
	return "client"
}

func (*clientCmd) Synopsis() string {
	return "connect to gRPC server"
}

func (*clientCmd) Usage() string {
	const use = `
Send ELF file parse request.

Usage:
	client [OPTION]... FILE

Flags:
`
	return use[1:]
}

func (cmd *clientCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.Addr, "addr", grpcAddr, "gRPC address to connect to")
}

func (cmd *clientCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	elfPath := f.Arg(0)

	// Connect to gRPC server.
	if err := connect(cmd.Addr, elfPath); err != nil {
		warn.Printf("connect failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// connect connects to the given gRPC address for incoming requests to parse ELF
// files.
func connect(addr, elfPath string) error {
	dbg.Printf("connecting to %q", addr)
	// Launch gRPC server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	// Send ELF parsing request.
	client := alvriket.NewELFParserClient(conn)
	ctx := context.Background()
	req := &alvriket.ParseELFRequest{
		ElfPath: elfPath,
	}
	reply, err := client.ParseELF(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("nsects:", reply.Nsects)
	return nil
}
