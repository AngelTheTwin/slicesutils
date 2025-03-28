package slicesutils

import (
	"fmt"
	"runtime"
)

// SafeExcecute executes a given function and recovers from any panic that occurs during its execution.
// It returns the output of the function and any error that occurred.
// If a panic occurs, it intercepts the panic and returns it as an error.
func SafeExcecute[T_out any](fn func() (T_out, error)) (output T_out, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	output, err = fn()
	return
}

// SafeExcecuteWithStackTrace executes a function that returns a value and an error,
// and ensures that any panic during the execution is recovered and converted into an error
// with a stack trace.
func SafeExcecuteWithStackTrace[T_out any](fn func() (T_out, error)) (output T_out, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			err = fmt.Errorf("panic: %v\nStack trace:\n%s", err, getErrWithStackTrace())
		}
	}()

	output, err = fn()
	return
}

func getErrWithStackTrace() string {
	buff := make([]byte, 4096)
	n := runtime.Stack(buff, false)
	stackTrace := string(buff[:n])
	return stackTrace
}
