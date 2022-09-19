package main

import (
	"context"
	"crypto/tls"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/openziti/sdk-golang/ziti/edge"
	cmap "github.com/orcaman/concurrent-map/v2"
	"io"
	"log"
	"net"
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

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	client := NewHttpClient(ctx, tlsConfig)

	resp, err := client.Get("http://" + serviceName)
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

func NewHttpClient(ctx ziti.Context, tlsConfig *tls.Config) *http.Client {
	return &http.Client{
		Transport: NewZitiTransport(ctx, tlsConfig),
	}
}

type ZitiTransport struct {
	http.Transport
	connByAddr cmap.ConcurrentMap[edge.Conn]
	Context    ziti.Context
	TlsConfig  *tls.Config
}

func NewZitiTransport(ctx ziti.Context, clientTlsConfig *tls.Config) *ZitiTransport {
	zitiTransport := &ZitiTransport{
		connByAddr: cmap.New[edge.Conn](),
		TlsConfig:  clientTlsConfig,
		Context:    ctx,
	}

	zitiTransport.Transport = http.Transport{
		DialContext:    zitiTransport.DialContext,
		DialTLSContext: zitiTransport.DialTLSContext,
	}

	return zitiTransport
}

func (transport *ZitiTransport) getConn(addr string) (edge.Conn, error) {
	var err error
	edgeConn := transport.connByAddr.Upsert(addr, nil, func(_ bool, existingConn edge.Conn, _ edge.Conn) edge.Conn {
		if existingConn == nil || existingConn.IsClosed() {
			var newConn edge.Conn

			cleanAddr := strings.Replace(addr, "http://", "", 1)
			cleanAddr = strings.Replace(cleanAddr, "https://", "", 1)
			cleanAddr = strings.Replace(cleanAddr, ":443", "", 1)
			cleanAddr = strings.Replace(cleanAddr, ":80", "", 1)
			newConn, err = transport.Context.Dial(cleanAddr)

			return newConn
		}

		return existingConn
	})

	return edgeConn, err
}

func (transport *ZitiTransport) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	edgeConn, err := transport.getConn(addr)

	return edgeConn, err
}

func (transport *ZitiTransport) DialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	edgeConn, err := transport.getConn(addr)

	if err != nil {
		return nil, err
	}
	tlsConn := tls.Client(edgeConn, transport.TlsConfig)

	if err := tlsConn.Handshake(); err != nil {
		if edgeConn != nil {
			_ = edgeConn.Close()
		}
		return nil, err
	}

	return edgeConn, err
}
