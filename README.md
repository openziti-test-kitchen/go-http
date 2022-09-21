# go-http

This project provides examples of how to use the Ziti GoLang SDK to both host and consume HTTP services over an
[OpenZiti overlay network](https://github.com/openziti/ziti). The Ziti SDK provides `ZitiTransport` which 
implements `http.Transport` and `edge.Conn` instances which can be used as `net.Listener` instances.

Links to the different projects and their example OpenZiti integrations can be found below.

# Servers

- [Standard Go Server](https://pkg.go.dev/net/http) - [example](./cmd/ziti-server-go/main.go)
- [Gin](https://github.com/gin-gonic/gin) - [example](./cmd/ziti-server-gin/main.go)
- [Goji](https://github.com/goji/goji) - [example](./cmd/ziti-server-goji/main.go)
- [Gorilla](https://github.com/gorilla/mux) - [example](./cmd/ziti-server-gorilla/main.go)

# Clients

- [Standard Go Client](https://pkg.go.dev/net/http) - [example](./cmd/ziti-client-go/main.go)
- [Resty](https://github.com/go-resty/resty) - [example](./cmd/ziti-client-resty/main.go)

# The Main Magic

GoLang's built in HTTP facilitates provide excellent methods for hooking into them. With the use of
the OpenZiti SDK it essentially boils down to the following patterns that work for the standard
GoLang HTTP client and server. These can then be adjusted to fit into any framework/library that uses
the GoLang HTTP packages.

The [OpenZiti GoLang SDK](https://github.com/openziti/sdk-golang) provides [`ZitiTransport`](https://github.com/openziti/sdk-golang/blob/main/http_transport.go),
which can be used as an `http.Transport`, and `edge.Listener`, that can be used as a `net.Listener`. `ZitiTransport`
can be used to create `http.Client` instances and `edge.Listener` can be used to with `http.Serve(listener,...)` calls.
The rest of the GoLang HTTP machinery handles all the HTTP interactions unknowingly over an OpenZiti network.

If you want to deep dive, the `ZitiTransport` definition can be found [here](https://github.com/openziti/sdk-golang/blob/main/http_transport.go)
and `edge.Listen()` can be found [here](https://github.com/openziti/sdk-golang/blob/main/ziti/ziti.go#L590).

## Client
Before:
```go
    client := http.DefaultClient
    resp, err := client.Get("http://" + args.ServiceName)
```

After:
```go
	client := sdk_golang.NewHttpClient(ctx, nil)
	resp, err := client.Get("http://" + args.ServiceName)
```

## Server
Before:
```go
	if err := http.Serve(listener, http.HandlerFunc(handler)); err != nil {
		log.Fatalf("serving failed: %v", err)
	}
```

After:
```go
	listener, err = ctx.Listen(args.ServiceName)

	if err := http.Serve(listener, http.HandlerFunc(handler)); err != nil {
		log.Fatalf("https serving failed: %v", err)
	}
```

# Example CLI Arguments

Each example uses the same [command line argument processing](./cmd/args.go). This processing takes in two
or four arguments that specify the Ziti Identity configuration file and OpenZiti service name. The two
additional arguments are paths to a x509 certificate and key in PEM format. If specified for a server,
the server will be hosted as an HTTPS service using the provided certificate and key files for the server's
identity. For a client, the x509 certificate and key will be used as the client certificate and key used to
initiate the TLS connection over the OpenZiti network.

`ziti-server-gin <serviceName> <identityConfig> [<certificate> <key>]`

However, HTTPS when working with OpenZiti is not necessary. See the next section!

# A Note on HTTPS

Hosting an HTTPS server over OpenZiti means that a TLS handshake will occur. A TLS handshake
requires that the server presents a certificate with a SAN IP or a SAN DNS entry that matches
the address the client used to access the service. For OpenZiti this means that a SAN DNS
that matches the OpenZiti service name must be present.

If the service will only be hosted over OpenZiti, HTTPS is an extra layer of security that can safely
be omitted. OpenZiti connections are inherently end-to-end encrypted and the data plane across
an OpenZiti network is additionally encrypted on each leg of transit. Additionally, the controller
has already verified all clients and hosts before they "dial" (connect) or "bind" (host).

# Setting Up The Examples

In order to run these examples, an OpenZiti network must be up and running. This includes a controller and router.
Additionally, a service, service host and client will need to be created. The host and client identities will need
policies to access and host the service. To set up an OpenZiti network, please see the 
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
5) Enroll your identities
    - `ziti edge enroll server.jwt` > creates `server.json`
    - `ziti edge enroll client.jwt` > creates `client.json`
6) Start an example
    - `ziti-server-gin myHttpService server.json`