package slicesutils

import (
	"fmt"
	"runtime"
)

// SafeExcecute executes a given function and recovers from any panic that occurs during its execution.
// It returns the output of the function and any error that occurred.
// If a panic occurs, it captures the stack trace and returns it as an error.
func SafeExcecute[T_out any](fn func() (T_out, error)) (output T_out, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = getErrWithStackTrace()
		}
	}()

	output, err = fn()
	return
}

func getErrWithStackTrace() error {
	buff := make([]byte, 4096)
	n := runtime.Stack(buff, false)
	stackTrace := string(buff[:n])
	err := fmt.Errorf("encountered panic: \nStack trace:\n%s", stackTrace)
	return err
}
