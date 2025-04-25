package trace

import (
	"time"
)

// ExecTimer create a debug log with the total time taken to complete a function execution.
//
// How to call this function? Should call at the start of the function
//
//	defer span.Finish()
//	defer ExecTimer("<log message>", span)()
func ExecTimer(message string, span ISpan) func() {
	start := time.Now()
	return func() {
		span.Debug(message, "executionTime", time.Since(start))
	}
}
