package main

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err := http.Serve(listener, r.Handler()); err != nil {
		log.Fatalf("https servering failed: %v", err)
	}
}
