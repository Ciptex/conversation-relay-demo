package trace

import (
	"log/slog"
	"os"
)

type ConsoleTrace struct {
	logLevel LogLevel
}

type ConsoleSpan struct {
	span     ISpan
	name     string
	logLevel LogLevel
	tags     []any
}

func newConsoleTrace(_ string, logLevel LogLevel) ITracer {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))

	slog.SetDefault(logger)

	trace := ConsoleTrace{
		logLevel: logLevel,
	}
	return &trace
}

func (t *ConsoleTrace) createSpanFrom(name string, spanRef map[string]string) ISpan {
	return ConsoleSpan{
		logLevel: t.logLevel,
		name:     name,
	}
}

func (t *ConsoleTrace) Start(name string) ISpan {
	return ConsoleSpan{
		logLevel: t.logLevel,
		name:     name,
	}
}

func (t *ConsoleTrace) startAsChildOf(name string, context interface{}) ISpan {
	return ConsoleSpan{
		logLevel: t.logLevel,
		name:     name,
	}
}

func (s ConsoleSpan) AsParent(name string) ISpan {
	return globalTracer.startAsChildOf(name, nil)
}

func (s ConsoleSpan) getContext() interface{} {
	return nil
}

func (s ConsoleSpan) SetTag(key, value string) ISpan {
	s.tags = append(s.tags, key, value)
	return s
}

func (s ConsoleSpan) GetSpanRef() map[string]string {
	tmp := map[string]string{}
	tmp["trace-id"] = "ctrace"
	return tmp
}

func (s ConsoleSpan) Info(message string, keyVals ...any) {
	if s.logLevel < 1 {
		return
	}
	go func() {
		args := s.msgArgs(keyVals...)
		slog.Info(message, args...)
	}()
}

func (s ConsoleSpan) Debug(message string, keyVals ...any) {
	if s.logLevel < 3 {
		return
	}
	go func() {
		args := s.msgArgs(keyVals...)
		slog.Debug(message, args...)
	}()
}

func (s ConsoleSpan) Dev(message string, keyVals ...any) {
	if s.logLevel < 4 {
		return
	}
	go func() {
		args := s.msgArgs(keyVals...)
		slog.Debug(message, args...)
	}()
}

func (s ConsoleSpan) Warn(message string, keyVals ...any) {
	if s.logLevel < 2 {
		return
	}
	go func() {
		args := s.msgArgs(keyVals...)
		slog.Warn(message, args...)
	}()
}

func (s ConsoleSpan) Error(message string, err error, keyVals ...any) {
	go func() {
		withName := keyVals
		if err != nil {
			withName = append(withName, "error", err.Error())
		}
		slog.Error(message, withName...)
	}()
}

func (s ConsoleSpan) Finish() {

}

func (s ConsoleSpan) msgArgs(keyVals ...any) []any {
	arr := keyVals
	arr = append(arr, s.tags[:]...)
	arr = append(arr, "name", s.name)
	return arr
}
