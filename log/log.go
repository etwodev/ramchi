package log

import (
	"context"
	"time"
)

// ctxKey is a private type used as a key for storing values in context.
// This prevents collisions with other context keys.
type ctxKey string

// loggerCtxKey is the key used to store the logger instance in the request context.
var LoggerCtxKey = ctxKey("logger")

type Logger interface {
	Debug() Entry
	Info() Entry
	Warn() Entry
	Error() Entry
	Fatal() Entry
}

type Entry interface {
	Str(key, value string) Entry
	Dur(key string, value time.Duration) Entry
	Int(key string, value int) Entry
	Bool(key string, value bool) Entry
	Msg(msg string)
	Err(error) Entry
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(LoggerCtxKey).(Logger); ok {
		return logger
	}
	return nil
}
