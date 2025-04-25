package trace

import (
	"errors"
)

var globalTracer ITracer

type ITracer interface {
	// Start create a new span instance with the passed in name
	Start(name string) ISpan

	startAsChildOf(name string, context interface{}) ISpan
	createSpanFrom(name string, spanRef map[string]string) ISpan
}

type ISpan interface {
	// SetTag set span level tags
	SetTag(key, value string) ISpan

	// Info log informational messages
	// 	message: string message to log
	// 	keyVals: comma seperated key value pair
	//	span.Info("an info message", "accSid", "ACXXXXX", "callSid", "CAXXXX")
	Info(message string, keyVals ...any)

	// Debug log debug messages
	// 	message: string message to log
	// 	keyVals: comma seperated key value pair
	//	span.Debug("a debug message", "accSid", "ACXXXXX", "callSid", "CAXXXX")
	Debug(message string, keyVals ...any)

	// Warn log warning messages
	// 	message: string message to log
	//	keyVals: comma seperated key value pair
	//	span.Warn("a warn message", "accSid", "ACXXXXX", "callSid", "CAXXXX")
	Warn(message string, keyVals ...any)

	// Dev log development messages
	// 	message: string message to log
	//	keyVals: comma seperated key value pair
	//	span.Dev("a dev message", "accSid", "ACXXXXX", "callSid", "CAXXXX")
	Dev(message string, keyVals ...any)

	// Error log error messages. Error messages will logged irrespective of log levels
	// 	message: string message to log
	//	err: error instance
	//	keyVals: comma seperated key value pair
	//	span.Error("an error occured", err, "accSid", "ACXXXXX", "callSid", "CAXXXX")
	Error(message string, err error, keyVals ...any)

	// AsParent creates a new child span from the current span
	AsParent(name string) ISpan

	// GetSpanRef get the span reference like span id and traceid etc
	GetSpanRef() map[string]string

	// Finish finish the span instance
	Finish()
}

type LogLevel int

const (
	OFF   LogLevel = 0
	INFO  LogLevel = 1
	WARN  LogLevel = 2
	DEBUG LogLevel = 3
	DEV   LogLevel = 4
	ALL   LogLevel = 5
)

type TraceType string

const (
	OpenTelemetry TraceType = "OPEN_TELIMETRY"
	Console       TraceType = "CONSOLE"
)

func CreateGlobalTracer(traceType TraceType, name string, loglevel LogLevel) (ITracer, error) {
	if traceType == OpenTelemetry {
		// globalTracer = newOpenTrace(name, loglevel)
	} else if traceType == Console {
		globalTracer = newConsoleTrace(name, loglevel)
	} else {
		return globalTracer, errors.New("invalid tracer type")
	}
	return globalTracer, nil
}

// GetGlobalTracer returns the global tracer created at startup
// Any modules can call this function and create a span to start logging
func GetGlobalTracer() ITracer {
	return globalTracer
}

// CreateChildSpanFrom create a child span from the passed spanRef
func CreateChildSpanFrom(spanRef map[string]string, name string) ISpan {
	if spanRef == nil {
		return globalTracer.Start(name)
	}
	return globalTracer.createSpanFrom(name, spanRef)
}
