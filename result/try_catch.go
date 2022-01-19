package result

import (
	"errors"
	"fmt"
)

// CallWithWrapPanic call func and wrap panic in error
func CallWithWrapPanic[T any](f func() T) (res Result[T]) {
	defer func() {
		if err := recover(); err != nil {
			var e error
			if _e, ok := err.(error); ok {
				e = _e
			} else {
				e = fmt.Errorf("panic: %v", e)
			}
			res = Err[T](e)
		}
	}()

	ret := f()
	return Ok(ret)
}

// Call call func, if panic by Result Unwrap, mute and instead of return Result.Err, otherwise propagate panic
func Call[T any](f func() T) (res Result[T]) {
	return CallExportPanicCustom(f, nil)
}

// CallExportPanicCustom call func,
// if panic by Result Unwrap, mute and instead of return Result.Err;
// else if dicider return bool, wrap panic error as Result.Err;
// else propagate panic.
func CallExportPanicCustom[T any](f func() T, dicider func(error) bool) (res Result[T]) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				if errors.Is(e, unwrapError) {
					rue := err.(*ResultUnwrapError[T])
					res = Err[T](rue.err)
					return
				}
				if dicider != nil && dicider(e) {
					res = Err[T](e)
					return
				}
			}
			panic(err)
		}
	}()

	ret := f()
	return Ok(ret)
}
