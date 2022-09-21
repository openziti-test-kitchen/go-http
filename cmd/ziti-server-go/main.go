package main

import (
	"crypto/tls"
	"github.com/openziti-test-kitchen/go-http/cmd"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"log"
	"net"
	"net/http"
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

	var listener net.Listener
	listener, err = ctx.Listen(args.ServiceName)

	if err != nil {
		log.Fatalf("could not bind service %s: %v", args.ServiceName, err)
	}

	if args.TlsConfig != nil {
		listener = tls.NewListener(listener, args.TlsConfig)
	}

	if err := http.Serve(listener, http.HandlerFunc(handler)); err != nil {
		log.Fatalf("https serving failed: %v", err)
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("Ziti HTTP!"))
}
