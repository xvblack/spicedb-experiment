package main

import (
	"context"
	"flag"
	"log"

	"github.com/authzed/authzed-go/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type secureMetadataCreds map[string]string

func (c secureMetadataCreds) RequireTransportSecurity() bool { return false }
func (c secureMetadataCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return c, nil
}

var server = flag.String("server", "localhost:50051", "server endpoint")
var preshared_key = flag.String("preshared_key", "somerandomkeyhere", "preshared key")

func main() {
	flag.Parse()
	client, err := authzed.NewClient(
		*server,
		grpc.WithPerRPCCredentials(secureMetadataCreds{"authorization": "Bearer " + *preshared_key}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInsecure(),
		// grpcutil.WithSystemCerts(grpcutil.SkipVerifyCA),
	)
	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	if flag.Arg(0) == "check" {
		run_check(client)
	} else {
		run_write(client)
	}
}
