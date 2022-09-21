package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/openziti-test-kitchen/go-http/cmd"
	"github.com/openziti/sdk-golang"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
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

	httpClient := sdk_golang.NewHttpClient(ctx, args.TlsConfig)

	restyClient := resty.NewWithClient(httpClient)

	resp, err := restyClient.R().Get("http://" + args.ServiceName)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	if resp != nil {
		log.Printf(string(resp.Body()))
	}
}
