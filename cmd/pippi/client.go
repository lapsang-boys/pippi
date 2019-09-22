package main

import (
	"context"
	"log"
	"time"

	stringspb "github.com/lapsang-boys/pippi/proto/strings"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func Strings(addr, binId string) ([]string, error) {
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

	rawStrings := []string{}
	for _, s := range reply.Strings {
		rawStrings = append(rawStrings, s.RawString)
	}
	return rawStrings, nil
}
