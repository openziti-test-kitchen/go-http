# go-http

This project provides examples of how to use the Ziti GoLang SDK to both host and consume HTTP services over an
[OpenZiti overlay network](https://github.com/openziti/ziti). The Ziti SDK provides `ZitiTransport` which 
implements `http.Transport` and `http.RountTripper` to intercept socket requests and instead provide 
`edge.Conn` instances which stand in as `net.Conn` interface implementors. This allows most GoLang HTTP clients 
and server libraries/framework to seamlessly work over an OpenZiti network.

Links to the different projects and their example OpenZiti integrations can be found below.

# Servers

- [Standard Go Server](https://pkg.go.dev/net/http) - [example](./cmd/ziti-server-go/main.go)
- [Gin](https://github.com/gin-gonic/gin) - [example](./cmd/ziti-server-gin/main.go)
- [Goji](https://github.com/goji/goji) - [example](./cmd/ziti-server-goji/main.go)
- [Gorilla](https://github.com/gorilla/mux) - [example](./cmd/ziti-server-gorilla/main.go)

# Clients

- [Standard Go Client](https://pkg.go.dev/net/http) - [example](./cmd/ziti-client-go/main.go)
- [Resty](https://github.com/go-resty/resty) - [example](./cmd/ziti-client-resty/main.go)


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

# Setting Up The Examples

In order to run these examples, an OpenZiti network must be up and running. This includes a controller and router.
Additionally, a service, service host and client will need to be created. The host and client identities will need
policies to access and host the service. To setup an OpenZiti network, please see the 
[quickstart guides](https://openziti.github.io/ziti/quickstarts/quickstart-overview.html).

You will need the [Ziti CLI](https://github.com/openziti/ziti/cmd/ziti) installed and on your path.

1) Login 
    - `ziti edge login "https://localhost:1280/edge/management/v1" -c $controllerCa -u $user -p $password`
2) Create the Service
    - `ziti create service myHttpService -a httpService`
3) Create the identities
    - `ziti create identity service httpServer -a httpServer -o server.jwt` > creates `server.jwt`
    - `ziti create identity user httpClient -a httpClient -o client.jwt` > creates `client.jwt`
4) Create policies
    - `ziti create service-policy httpServers Bind --identity-roles #httpServer --service-roles #httpService`
    - `ziti create service-policy httpClients Dial --identity-roles #httpClient --service-roles #httpService`
6) Enroll your identities
    - `ziti edge enroll server.jwt` > creates `server.json`
    - `ziti edge enroll client.jwt` > creates `client.json`
7) Start an example
    - `ziti-server-gin myHttpService server.json`
