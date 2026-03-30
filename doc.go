// Package stacklog provides lightweight error tracing and colored console logging.
//
// Quick usage:
//
//	err := doWork()
//	if err != nil {
//	    return stacklog.Trace(err) // adds [file:line] and preserves nesting
//	}
//
// For structured app logs, create an APIPrint once per service and call Info/Error.
// The functions are documented with brief examples so editors show helpful hover text.
package stacklog
