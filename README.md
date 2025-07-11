# ramchi

`ramchi` is an extension to the [go-chi](https://github.com/go-chi/chi) HTTP router designed for rapid and modular development of REST APIs. `ramchi` emphasizes developer experience.

## Features

- Modular router and middleware loading system
- Feature flag support via experimental toggles for safe progressive rollout
- Easy-to-configure TLS/HTTPS support
- Graceful shutdown and OS signal handling for clean service termination
- Built-in, extensible helper packages for requests, responses, crypto, email, and more
- Automatic configuration file generation and hot loading
- Integrated structured logging powered by [zerolog](https://github.com/rs/zerolog)
- Toggle middleware globally via configuration, reducing boilerplate

## Installation

Use Go Modules to install:

```bash
go get -u github.com/etwodev/ramchi
```

## Getting Started

Below is a minimal example demonstrating how to create and start a `ramchi` server with a modular router and endpoint registration.

```go
package main

import (
	"encoding/json"
	"net/http"

	"github.com/etwodev/ramchi"
	"github.com/etwodev/ramchi/router"
)

func main() {
	s := ramchi.New()
	s.LoadRouter(Routers())
	s.Start()
}

// Define routers with their prefixes and routes
func Routers() []router.Router {
	return []router.Router{
		router.NewRouter("example", Routes(), true, nil),
	}
}

// Define individual routes for the router
func Routes() []router.Route {
	return []router.Route{
		router.NewGetRoute("/demo", true, false, ExampleGetHandler, nil),
	}
}

// ExampleGetHandler handles GET /example/demo requests
func ExampleGetHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(map[string]string{"success": "ping"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(res); err != nil {
		// Handle write error (in production, use proper error logging)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
```

> **Note:** On the first run, `ramchi` auto-generates a default `ramchi.config.json` file in your working directory, which you can customize as needed.

## Configuration

The behavior of the server and feature toggling is controlled by the `ramchi.config.json` file.

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
| `tlsCertFile`          | string    | Path to TLS certificate file (required if `enableTLS` is `true`)            | `""`        |
| `tlsKeyFile`           | string    | Path to TLS key file (required if `enableTLS` is `true`)                    | `""`        |
| `readTimeout`          | int       | Max seconds allowed to read incoming requests                               | `15`        |
| `writeTimeout`         | int       | Max seconds allowed to write responses                                      | `15`        |
| `idleTimeout`          | int       | Max seconds to keep idle HTTP connections open                              | `60`        |
| `maxHeaderBytes`       | int       | Maximum size in bytes for HTTP headers                                      | `1048576`   |
| `shutdownTimeout`      | int       | Timeout in seconds for graceful server shutdown                             | `15`        |
| `enableCORS`           | bool      | Automatically enables Cross-Origin Resource Sharing (CORS) middleware       | `false`     |
| `allowedOrigins`       | \[]string | List of allowed origins for CORS (e.g., `["*"]`, `["https://example.com"]`) | `["*"]`     |
| `enableRequestLogging` | bool      | Enables HTTP request logging middleware globally                            | `false`     |

## Togglable Middleware

Enable built-in middleware globally using config flags without manual registration:

| Middleware          | Config Flag            | Description                                      |
| ------------------- | ---------------------- | ------------------------------------------------ |
| **CORS**            | `enableCORS`           | Adds permissive or restricted CORS headers       |
| **Request Logging** | `enableRequestLogging` | Logs incoming HTTP requests with structured logs |

> For fine-grained control (middleware order, conditional logic), register middleware manually using `LoadMiddleware()`.

## Accessing Configuration in Code

You can access configuration values in your application via the `config` package:

```go
import c "github.com/etwodev/ramchi/v2/config"

port := c.Port()           // e.g. "7000"
addr := c.Address()        // e.g. "0.0.0.0"

if c.Experimental() {
	// Enable or disable experimental features accordingly
}

if c.EnableTLS() {
	cert := c.TLSCertFile()
	key := c.TLSKeyFile()
	// Use cert and key for HTTPS server setup
}
```

## Logging

`ramchi` integrates [zerolog](https://github.com/rs/zerolog) for structured, leveled logging:

* Log levels controlled via `logLevel` config (`debug`, `info`, `warn`, `error`, `fatal`, `disabled`)
* Console-friendly output by default, but pluggable to other loggers if needed
* Logs contextual information: server name, HTTP method, route, middleware names, error stack traces
* Logs graceful shutdown steps, warnings, and fatal errors

## Middleware & Routing Best Practices

* Organize routes modularly using `router.Router` instances grouped by prefixes.
* Respect feature flags by setting the `Experimental` flag on routes/middleware.
* Use middleware chaining to add cross-cutting concerns like authentication, CORS, logging.
* Use the status flag to disable routes/middleware temporarily without deleting code.

## TLS Support

To serve over HTTPS:

1. Set `"enableTLS": true` in `ramchi.config.json`.
2. Provide valid paths to `"tlsCertFile"` and `"tlsKeyFile"` for your SSL certificate and private key.
3. Restart your server.

`ramchi` will handle HTTPS setup automatically.

## Extending ramchi with Helpers

`ramchi` includes utility helper packages to accelerate development:

| Helper Package        | Description                                                |
| --------------------- | ---------------------------------------------------------- |
| `helpers/request.go`  | Utilities for HTTP requests (client IP extraction, params) |
| `helpers/response.go` | JSON encoding, error handling, standardized responses      |
| `helpers/crypto.go`   | Cryptographic helpers: hashing, encryption utilities       |
| `helpers/email.go`    | Email templating and sending utilities                     |
| `helpers/strings.go`  | String manipulation (padding, truncation, sanitization)    |
| `helpers/encoding.go` | Encoding helpers (hex, base64 conversions)                 |

Feel free to extend or create your own helper packages and contribute back.

## Contribution Guidelines

Contributions, feature requests, and bug reports are welcome! Please:

* Fork the repository
* Create a feature branch (`git checkout -b feature-name`)
* Write tests for your changes
* Submit a Pull Request describing your improvements
* Open issues for discussion before implementing breaking changes

## Example Advanced Usage

Here's an example of grouping multiple routers and adding middleware.
Of course, in reality, you would put your routes in a separate package for brevity:

```go
package main

import (
	"net/http"

	"github.com/example/project/middleware/logger"

	"github.com/etwodev/ramchi"
	"github.com/etwodev/ramchi/router"
)

func main() {
	s := ramchi.New()
	s.LoadMiddleware(Middlewares())
	s.LoadRouter(Routers())
	s.Start()
}

func Routers() []router.Router {
	return []router.Router{
		router.NewRouter("api/v1", apiV1Routes(), true, nil),
		router.NewRouter("admin", adminRoutes(), true, nil),
	}
}

func Middlewares() []middleware.Middleware {
	return []middleware.Middleware{
		middleware.NewMiddleware(logger.Middleware(), "logger", true, false),
	}
}

func apiV1Routes() []router.Route {
	return []router.Route{
		router.NewGetRoute("/users", true, false, usersHandler, nil),
		router.NewPostRoute("/users", true, false, createUserHandler, nil),
	}
}

func adminRoutes() []router.Route {
	return []router.Route{
		router.NewGetRoute("/dashboard", true, false, adminDashboardHandler, nil),
	}
}
```

## Contact and Support

For questions, discussions, or support, please open an issue.

## License

MIT License © Etwodev

