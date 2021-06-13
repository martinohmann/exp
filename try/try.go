// Package try provides an experimental approach to error handling which uses
// panics internally to control execution flow. It can be used to avoid
//
//   if err != nil {
//       return err
//   }
//
// checks by replacing them with
//
//   try.Check(err)
//
// which will abort execution when err is non-nil.
package try

import "fmt"

type tryError struct {
	error
}

func (e *tryError) Unwrap() error {
	return e.error
}

// Check should be called inside the closure passed to Run to handle errors.
// If err is non-nil it panics with a custom error that is handled and
// recovered by Run. If Check is called with a non-nil error outside of Run it
// panics.
func Check(err error) {
	if err != nil {
		panic(&tryError{err})
	}
}

// Checkf behaves like Check but can be used to wrap err with additional
// context. Given the following call:
//
//   try.Checkf(err, "error context: %s", someString)
//
// For a non-nil error the call aborts the execution flow with an error equivalent to:
//
//   fmt.Errorf("error context: %s: %w", someString, err)
//
// See documentation on Check for more information.
func Checkf(err error, format string, args ...interface{}) {
	if err != nil {
		panic(&tryError{fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)})
	}
}

// Run executes fn and aborts the execution flow if Check or Checkf are invoked
// inside fn with a non-nil error. If a panic occurs inside fn, Run panics as
// well. Returns the first non-nil error passed to an invocation of Check or
// Checkf.
//
//   err := try.Run(func() {
//       try.Check(doStuff())
//
//       err := doSomethingElse()
//       try.Checkf(err, "failed to do something else")
//   })
//   if err != nil {
//       // handle error
//   }
//
func Run(fn func()) (err error) {
	defer Catch(&err)
	fn()
	return
}

// Catch recovers from errors if the execution flow is aborted when Check or
// Checkf are invoked with a non-nil error. Sets err to the value of the first
// non-nil error passed to Check or Checkf.
//
//   func fn() (err error) {
//       defer try.Catch(&err)
//
//       try.Checkf(doSomething(), "failed to do something")
//       try.Check(doSomethingElse())
//
//       return err
//   }
//
func Catch(err *error) {
	switch e := recover().(type) {
	case nil:
		return
	case *tryError:
		*err = e.Unwrap()
	default:
		panic(e)
	}
}
