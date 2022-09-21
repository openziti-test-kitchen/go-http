# go-http

This project provides examples of how to use the Ziti GoLang SDK to both host and consume HTTP services over an
[OpenZiti overlay network](https://github.com/openziti/ziti). The Ziti SDK provides `ZitiTransport` which 
implements `http.Transport` and `http.RountTripper` to intercept socket requests and instead provide 
`edge.Conn` instances which stand in as `net.Conn` interface implementors. This allows most GoLang HTTP clients 
and server libraries/framework to seamlessly work over an OpenZiti network.

Links to the different projects and their example OpenZiti integrations can be found below.

# Servers

- [Standard Go Server](https://pkg.go.dev/net/http) - [example](./cmd/ziti-server-go)
- [Gin](https://github.com/gin-gonic/gin) - [example](./cmd/ziti-server-gin)
- [Goji](https://github.com/goji/goji) - [example](./cmd/ziti-server-goji)
- [Gorilla](https://github.com/gorilla/mux) - [example](./cmd/ziti-server-gorilla)

# Clients

- [Standard Go Client](https://pkg.go.dev/net/http) - [example](./cmd/ziti-client-go)
- [Resty](https://github.com/go-resty/resty) - [example](./cmd/ziti-client-resty)


# Explanation of Examples

# CLI Arguments

Each example uses the same [command line argument processing](./cmd/args.go). This processing takes in two
or four arguments that specify the Ziti Identity configuration file and OpenZiti service name. The two
additional arguments are paths to a x509 certificate and key in PEM format. If specified for a server,
the server will be hosted as an HTTPS service using the provided certificate and key files for the server's
identity. For a client, the x509 certificate and key will be used as the client certificate and key used to
initiate the TLS connection over the OpenZiti network.

`ziti-server-gin <serviceName> <identityConfig> [<certificate> <key>]`

# ZitiTransport

The [OpenZiti GoLang SDK](https://github.com/openziti/sdk-golang) provides a `ZitiTransport` which can be used
as an `http.Transport`. This effectively reduces all examples to providing an `http.Client` that uses a
`ZitiTransport` instance that implements `http.RountTripper`. The rest of the GoLang HTTP machinery handles all
the HTTP interactions unknowing over an OpenZiti network.