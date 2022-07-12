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

func main() {
	flag.Parse()
	if flag.Arg(0) == "check" {
		run_check()
	} else {
		run_write()
	}
}
