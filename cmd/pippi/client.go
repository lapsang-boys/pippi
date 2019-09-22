package main

import (
	"context"
	"log"
	"time"

	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	stringspb "github.com/lapsang-boys/pippi/proto/strings"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func Strings(addr, binId string) ([]*stringspb.StringInfo, error) {
	// Connect to gRPC server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := stringspb.NewStringsExtractorClient(conn)
	ctx := context.Background()

	req := &stringspb.StringsRequest{
		Id: binId,
	}
	now := time.Now()
	reply, err := client.ExtractStrings(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println(time.Since(now))

	return reply.Strings, nil
}

func Sections(addr, binId string) (*binpb.File, error) {
	// Connect to gRPC server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := binpb.NewBinaryParserClient(conn)
	ctx := context.Background()

	req := &binpb.ParseBinaryRequest{
		BinId: binId,
	}
	now := time.Now()
	file, err := client.ParseBinary(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println(time.Since(now))

	return file, nil
}

func Disassembly(addr, binId string) (*disasmpb.DisassembleReply, error) {
	// Connect to gRPC server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := disasmpb.NewDisassemblerClient(conn)
	ctx := context.Background()

	req := &disasmpb.DisassembleRequest{
		BinId: binId,
	}
	now := time.Now()
	reply, err := client.Disassemble(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println(time.Since(now))

	return reply, nil
}
