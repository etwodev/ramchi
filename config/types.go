package config

// Config holds all the configurable parameters for the application.
// It is serialized/deserialized from JSON config file.
type Config struct {
	Port                 string   `json:"port"`                 // the port to use
	Address              string   `json:"address"`              // the address to use
	Experimental         bool     `json:"experimental"`         // whether or not to enable experimental middleware/endpoints
	ReadTimeout          int      `json:"readTimeout"`          // seconds
	WriteTimeout         int      `json:"writeTimeout"`         // seconds
	IdleTimeout          int      `json:"idleTimeout"`          // seconds
	LogLevel             string   `json:"logLevel"`             // e.g. "debug", "info", "disabled"
	MaxHeaderBytes       int      `json:"maxHeaderBytes"`       // the maximum number of bytes in a request header
	EnableTLS            bool     `json:"enableTLS"`            // whether or not TLS should be enabled
	TLSCertFile          string   `json:"tlsCertFile"`          // if TLS is in use, the file path for the certificate
	TLSKeyFile           string   `json:"tlsKeyFile"`           // if TLS is in use, the file path for the key
	ShutdownTimeout      int      `json:"shutdownTimeout"`      // graceful shutdown timeout seconds
	EnableCORS           bool     `json:"enableCORS"`           // whether the CORS middleware should be enabled
	AllowedOrigins       []string `json:"allowedOrigins"`       // the allowed origins for CORS
	EnableRequestLogging bool     `json:"enableRequestLogging"` // whether request logging middleware should be enabled
}

// Port returns the configured server port.
func Port() string { return c.Port }

// Address returns the configured server address.
func Address() string { return c.Address }

// Experimental returns whether experimental features are enabled.
func Experimental() bool { return c.Experimental }

// ReadTimeout returns the server read timeout duration in seconds.
func ReadTimeout() int { return c.ReadTimeout }

// WriteTimeout returns the server write timeout duration in seconds.
func WriteTimeout() int { return c.WriteTimeout }

// IdleTimeout returns the server idle timeout duration in seconds.
func IdleTimeout() int { return c.IdleTimeout }

// LogLevel returns the configured logging level.
func LogLevel() string { return c.LogLevel }

// MaxHeaderBytes returns the maximum size of request headers in bytes.
func MaxHeaderBytes() int { return c.MaxHeaderBytes }

// EnableTLS indicates if TLS is enabled.
func EnableTLS() bool { return c.EnableTLS }

// TLSCertFile returns the path to the TLS certificate file.
func TLSCertFile() string { return c.TLSCertFile }

// TLSKeyFile returns the path to the TLS key file.
func TLSKeyFile() string { return c.TLSKeyFile }

// ShutdownTimeout returns the graceful shutdown timeout duration in seconds.
func ShutdownTimeout() int { return c.ShutdownTimeout }

// EnableCORS returns true if CORS support is enabled.
func EnableCORS() bool { return c.EnableCORS }

// AllowedOrigins returns the list of allowed CORS origins.
func AllowedOrigins() []string { return c.AllowedOrigins }

// EnableRequestLogging indicates if request logging is enabled.
func EnableRequestLogging() bool { return c.EnableRequestLogging }
