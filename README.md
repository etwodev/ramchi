# ramchi

`ramchi` is an extension to the [chi](https://github.com/go-chi/chi) HTTP router designed for rapid and modular development of web applications. `ramchi` focuses on developer experience while ensuring your website remains fast and responsive.

---

## Features

- Modular router and middleware loading
- Support for feature flagging via experimental toggles
- Unified backend and frontend serving capabilities
- Zero-configuration TLS support
- Structured, leveled logging powered by `zerolog`
- Graceful shutdown and signal handling
- Extensible helpers for requests, responses, crypto, email, and more

---

## Installation

```bash
go get -u github.com/etwodev/ramchi
```

---

## Getting Started

ramchi allows easy, modular registration of endpoints through grouping.

Create a new server instance and start it:

```go
package main

import (
  "github.com/etwodev/ramchi"
  "encoding/json"
	"net/http"
)

func main() {
	s := ramchi.New()
	s.LoadRouter(Routers())
	s.Start()
}

func Routers() []router.Router {
	return []router.Router{
		router.NewRouter("example", Routes(), true),
	}
}

func Routes() []router.Route {
	return []router.Route{
		router.NewGetRoute("/demo", true, false, ExampleGetHandler),
	}
}

// This route will be a GET endpoint registered at /example/demo
func ExampleGetHandler(w http.ResponseWriter, r *http.Request) {
  res, _ := json.Marshal(map[string]string{"success": "ping"})
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(201)
  if _, err := w.Write(res); err != nil {
    t.Fatal(err)
  }
}
```

On the first run, `ramchi` will automatically generate a default `ramchi.config.json` file in your working directory.

---

## Configuration

The `ramchi.config.json` file controls server behavior and feature toggling.

### Default Configuration Example

```json
{
  "port": "7000",
  "address": "0.0.0.0",
  "experimental": false,
  "logLevel": "info",
  "enableTLS": false,
  "tlsCertFile": "",
  "tlsKeyFile": "",
  "readTimeout": 15,
  "writeTimeout": 15,
  "idleTimeout": 60,
  "maxHeaderBytes": 1048576,
  "shutdownTimeout": 15,
  "enableCORS": false,
  "allowedOrigins": ["*"],
  "enableRequestLogging": false
}
```

### Configuration Fields

| Field                  | Type      | Description                                                                 | Default     |
| ---------------------- | --------- | --------------------------------------------------------------------------- | ----------- |
| `port`                 | string    | TCP port the server listens on                                              | `"7000"`    |
| `address`              | string    | IP address to bind to                                                       | `"0.0.0.0"` |
| `experimental`         | bool      | Enables or disables experimental feature flags                              | `false`     |
| `logLevel`             | string    | Log verbosity level (`debug`, `info`, `warn`, `error`, `fatal`, `disabled`) | `"info"`    |
| `enableTLS`            | bool      | Enable HTTPS by providing TLS certificate and key                           | `false`     |
| `tlsCertFile`          | string    | Path to TLS certificate file (required if `enableTLS` is true)              | `""`        |
| `tlsKeyFile`           | string    | Path to TLS key file (required if `enableTLS` is true)                      | `""`        |
| `readTimeout`          | int       | Maximum duration (in seconds) for reading the request                       | `15`        |
| `writeTimeout`         | int       | Maximum duration (in seconds) before timing out response writes             | `15`        |
| `idleTimeout`          | int       | Maximum duration (in seconds) to keep idle connections open                 | `60`        |
| `maxHeaderBytes`       | int       | Maximum size of request headers in bytes                                    | `1048576`   |
| `shutdownTimeout`      | int       | Time (in seconds) allowed for graceful shutdown                             | `15`        |
| `enableCORS`           | bool      | Automatically enables CORS middleware                                       | `false`     |
| `allowedOrigins`       | \[]string | List of allowed CORS origins (e.g., `["*"]`, `["https://example.com"]`)     | `["*"]`     |
| `enableRequestLogging` | bool      | Automatically enables HTTP request logging middleware                       | `false`     |

---

## Togglable Middleware

You can enable built-in middleware through the config file without registering them manually.

### Available Middleware

| Middleware          | Config Flag            | Description                                                     |
| ------------------- | ---------------------- | --------------------------------------------------------------- |
| **CORS**            | `enableCORS`           | Adds a permissive or origin-restricted CORS layer               |
| **Request Logging** | `enableRequestLogging` | Logs all incoming HTTP requests using structured logging format |

These are injected globally before any custom middleware or routes.

If you require more control (e.g., middleware ordering or conditional logic), you can still register them manually through the `LoadMiddleware()` method.

---

## Using the Configuration in Code

Your application can access config values via the `config` package accessor functions:

```go
import c "github.com/Etwodev/ramchi/v2/config"

port := c.Port()                 // e.g. "7000"
address := c.Address()           // e.g. "0.0.0.0"
if c.Experimental() {
    // Enable experimental features
}
level := c.LogLevel()            // e.g. "debug"
if c.EnableTLS() {
    cert := c.TLSCertFile()
    key := c.TLSKeyFile()
    // Use cert and key to start HTTPS server
}
```

The server internally uses these config values to set up logging, timeouts, TLS, and feature flags.

---

## Logging

`ramchi` integrates [zerolog](https://github.com/rs/zerolog) for structured, leveled logging with console-friendly output by default. However, logging can be replaced. If you would like a specific package to be supported, please raise an issue.

* Log verbosity is controlled by the `logLevel` config (e.g., `debug`, `info`, `disabled`).
* Logs include contextual fields such as the server group, function names, HTTP method, route path, and middleware names.
* Graceful shutdown logs warnings and fatal errors as appropriate.

---

## Middleware & Routing

* Load your middlewares and routers modularly before starting the server.
* `ramchi` respects middleware and route `Experimental` flags based on your config.
* Routes and middleware with disabled status or mismatched experimental flags are skipped.

---

## TLS Support

Set `enableTLS` to `true` and provide valid paths to `tlsCertFile` and `tlsKeyFile` in your config to serve HTTPS.

---

## Extending Helpers

`ramchi` ships with helper packages for common tasks:

* `helpers/request.go`: HTTP request utilities (e.g., extracting IP, URL params)
* `helpers/response.go`: Response helpers for JSON encoding, error handling
* `helpers/crypto.go`: Crypto utilities (hashing, encryption helpers)
* `helpers/email.go`: Email sending and templating helpers
* `helpers/strings.go`: String manipulation utilities (e.g., truncation, padding, sanitization)
* `helpers/encoding.go`: Encoding utilities (e.g., toHex, toBase64)

You are encouraged to extend these helper packages or create your own.

---

## Contributing

Contributions and suggestions are welcome. Please open issues or pull requests on the [GitHub repository](https://github.com/etwodev/ramchi).

---

## License

MIT License © Etwodev
