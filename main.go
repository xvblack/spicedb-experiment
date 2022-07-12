package main

import (
	"context"
	"flag"
)

type secureMetadataCreds map[string]string

func (c secureMetadataCreds) RequireTransportSecurity() bool { return false }
func (c secureMetadataCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return c, nil
}

var server = flag.String("server", "localhost:50051", "server endpoint")

func main() {
	flag.Parse()
	if flag.Arg(0) == "check" {
		run_check(*server)
	} else {
		run_write(*server)
	}
}
