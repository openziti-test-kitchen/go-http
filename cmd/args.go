package cmd

import (
	"crypto/tls"
	"log"
	"os"
	"strings"
)

type Arguments struct {
	ServiceName string
	ConfigPath  string
	CertPath    string
	KeyPath     string
	TlsConfig   *tls.Config
}

func ParseArgs() Arguments {
	result := Arguments{}

	args := os.Args[1:]
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		log.Fatalf("expected at least 2 arguments or 4 got %d", len(args))
	}

	result.ServiceName = strings.TrimSpace(args[0])

	if len(args) > 1 {
		result.ConfigPath = strings.TrimSpace(args[1])
	}

	if result.ConfigPath == "" {
		log.Fatalf("expected config file as second argument")
	}

	if len(args) == 4 {
		result.CertPath = args[2]
		result.KeyPath = args[3]

		tlsCert, err := tls.LoadX509KeyPair(result.CertPath, result.KeyPath)

		if err != nil {
			log.Fatalf("could not load certificate/key pair: %v", err)
		}

		result.TlsConfig = &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		}
	}

	return result
}
