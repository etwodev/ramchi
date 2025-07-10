package log

import "time"

type NoOpLogger struct{}

func (l *NoOpLogger) Debug() Entry { return &noopEntry{} }
func (l *NoOpLogger) Info() Entry  { return &noopEntry{} }
func (l *NoOpLogger) Warn() Entry  { return &noopEntry{} }
func (l *NoOpLogger) Error() Entry { return &noopEntry{} }
func (l *NoOpLogger) Fatal() Entry { return &noopEntry{} }

type noopEntry struct{}

func (n *noopEntry) Str(string, string) Entry        { return n }
func (n *noopEntry) Dur(string, time.Duration) Entry { return n }
func (n *noopEntry) Int(string, int) Entry           { return n }
func (n *noopEntry) Bool(string, bool) Entry         { return n }
func (n *noopEntry) Err(error) Entry                 { return n }
func (n *noopEntry) Msg(string)                      {}
