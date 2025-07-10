package config

type Config struct {
	Port                 string   `json:"port"`
	Address              string   `json:"address"`
	Experimental         bool     `json:"experimental"`
	ReadTimeout          int      `json:"readTimeout"`  // in seconds
	WriteTimeout         int      `json:"writeTimeout"` // in seconds
	IdleTimeout          int      `json:"idleTimeout"`  // in seconds
	LogLevel             string   `json:"logLevel"`     // e.g. "debug", "info", "disabled"
	MaxHeaderBytes       int      `json:"maxHeaderBytes"`
	EnableTLS            bool     `json:"enableTLS"`
	TLSCertFile          string   `json:"tlsCertFile"`
	TLSKeyFile           string   `json:"tlsKeyFile"`
	ShutdownTimeout      int      `json:"shutdownTimeout"` // graceful shutdown timeout seconds
	EnableCORS           bool     `json:"enableCORS"`
	AllowedOrigins       []string `json:"allowedOrigins"`
	EnableRequestLogging bool     `json:"enableRequestLogging"`
}

func Port() string               { return c.Port }
func Address() string            { return c.Address }
func Experimental() bool         { return c.Experimental }
func ReadTimeout() int           { return c.ReadTimeout }
func WriteTimeout() int          { return c.WriteTimeout }
func IdleTimeout() int           { return c.IdleTimeout }
func LogLevel() string           { return c.LogLevel }
func MaxHeaderBytes() int        { return c.MaxHeaderBytes }
func EnableTLS() bool            { return c.EnableTLS }
func TLSCertFile() string        { return c.TLSCertFile }
func TLSKeyFile() string         { return c.TLSKeyFile }
func ShutdownTimeout() int       { return c.ShutdownTimeout }
func EnableCORS() bool           { return c.EnableCORS }
func AllowedOrigins() []string   { return c.AllowedOrigins }
func EnableRequestLogging() bool { return c.EnableRequestLogging }
