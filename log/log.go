package log

import "time"

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
