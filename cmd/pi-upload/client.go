package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/google/subcommands"
	uploadpb "github.com/lapsang-boys/pippi/proto/upload"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// clientCmd is the command to upload a binary file.
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

func newRequest(binPath string) (*uploadpb.UploadRequest, error) {
	buf, err := ioutil.ReadFile(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rawHash := sha256.Sum256(buf)
	hash := hex.EncodeToString(rawHash[:])

	req := &uploadpb.UploadRequest{
		Filename: binPath,
		Hash:     hash,
		Content:  buf,
	}

	return req, nil
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
	client := uploadpb.NewUploadClient(conn)
	ctx := context.Background()

	req, err := newRequest(binPath)
	if err != nil {
		return errors.WithStack(err)
	}

	reply, err := client.Upload(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println(reply.Id)

	return nil
}