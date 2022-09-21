package main

import (
	"github.com/openziti-test-kitchen/go-http/cmd"
	sdk_golang "github.com/openziti/sdk-golang"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"io"
	"log"
)

func main() {
	args := cmd.ParseArgs()

	identityConfig, err := config.NewFromFile(args.ConfigPath)

	if err != nil {
		log.Fatalf("could not load configuration file: %v", err)
	}

	ctx := ziti.NewContextWithConfig(identityConfig)

	if err = ctx.Authenticate(); err != nil {
		log.Fatalf("could not authenticate: %v", err)
	}

	client := sdk_golang.NewHttpClient(ctx, args.TlsConfig)

	resp, err := client.Get("http://" + args.ServiceName)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	if resp != nil {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Error reading body: %v", err)
		}

		log.Printf(string(body))
	}
}
