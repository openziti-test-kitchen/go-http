package main

import (
	"crypto/tls"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		log.Fatalf("expected at least 2 arguments or 4 got %d", len(args))
	}

	serviceName := strings.TrimSpace(args[0])
	configFile := ""

	if len(args) > 1 {
		configFile = strings.TrimSpace(args[1])
	}

	if configFile == "" {
		log.Fatalf("expected config file")
	}

	identityConfig, err := config.NewFromFile(configFile)

	if err != nil {
		log.Fatalf("could not load identity configuration file: %v", err)
	}

	ctx := ziti.NewContextWithConfig(identityConfig)

	if err = ctx.Authenticate(); err != nil {
		log.Fatalf("could not authenticate with controller: %v", err)
	}

	listener, err := ctx.Listen(serviceName)

	if err != nil {
		log.Fatalf("could not bind service %s: %v", serviceName, err)
	}

	if len(args) == 4 {
		log.Printf("running HTTPS server, expecting FQDN to be in the certificate file SANs")

		tlsCert, err := tls.LoadX509KeyPair(args[2], args[3])

		if err != nil {
			log.Fatalf("could not load certificate/key pair: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		}

		tlsListener := tls.NewListener(listener, tlsConfig)

		if err := http.Serve(tlsListener, http.HandlerFunc(httpHandler)); err != nil {
			log.Fatalf("https servering failed: %v", err)
		}
	} else {
		if err := http.Serve(listener, http.HandlerFunc(httpHandler)); err != nil {
			log.Fatalf("http servering failed: %v", err)
		}
	}
}

func httpHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("Ziti HTTP!"))
}
