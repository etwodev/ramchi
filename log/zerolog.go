package log

import (
	"time"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	z zerolog.Logger
}

func NewZeroLogger(z zerolog.Logger) *ZeroLogger {
	return &ZeroLogger{z}
}

func (zl *ZeroLogger) Debug() Entry { return &zeroEntry{zl.z.Debug()} }
func (zl *ZeroLogger) Info() Entry  { return &zeroEntry{zl.z.Info()} }
func (zl *ZeroLogger) Warn() Entry  { return &zeroEntry{zl.z.Warn()} }
func (zl *ZeroLogger) Error() Entry { return &zeroEntry{zl.z.Error()} }
func (zl *ZeroLogger) Fatal() Entry { return &zeroEntry{zl.z.Fatal()} }

type zeroEntry struct {
	e *zerolog.Event
}

func (z *zeroEntry) Str(k, v string) Entry               { z.e.Str(k, v); return z }
func (z *zeroEntry) Dur(k string, v time.Duration) Entry { z.e.Dur(k, v); return z }
func (z *zeroEntry) Int(k string, v int) Entry           { z.e.Int(k, v); return z }
func (z *zeroEntry) Bool(k string, v bool) Entry         { z.e.Bool(k, v); return z }
func (z *zeroEntry) Err(e error) Entry                   { z.e.Err(e); return z }
func (z *zeroEntry) Msg(m string)                        { z.e.Msg(m) }
