package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	binpb "github.com/lapsang-boys/pippi/proto/bin"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// clientCmd is the command to connect to a gRPC server to send parse binary
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
Send binary file parse request.

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
	binPath := f.Arg(0)

	// Connect to gRPC server.
	if err := connect(cmd.Addr, binPath); err != nil {
		warn.Printf("connect failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// connect connects to the given gRPC address for incoming requests to parse
// binary files.
func connect(addr, binPath string) error {
	dbg.Printf("connecting to %q", addr)
	// Launch gRPC server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := binpb.NewBinaryParserClient(conn)
	ctx := context.Background()
	req := &binpb.ParseBinaryRequest{
		BinPath: binPath,
	}
	reply, err := client.ParseBinary(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("nsects:", reply.Nsects)
	return nil
}
