package main

import (
	"crypto/tls"
	"fmt"
	"github.com/openziti-test-kitchen/go-http/cmd"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	goji "goji.io"
	"goji.io/pat"
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

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/hello/:name"), hello)

	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("https servering failed: %v", err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	fmt.Fprintf(w, "Hello, %s!", name)
}
