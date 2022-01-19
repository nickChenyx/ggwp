package result

import (
	"errors"
	"fmt"
)

type Result[T any] struct {
	result T
	err    error
}

var (
	unwrapError = errors.New("Result[T] unwrap error")
)

type ResultUnwrapError[T any] struct {
	result T
	err    error
}

func (r *ResultUnwrapError[T]) Error() string {
	return fmt.Sprintf("Result[%T].IsErr() true", r.result)
}

func (r *ResultUnwrapError[T]) Unwrap() error { return r }

func (r *ResultUnwrapError[T]) Is(err error) bool { return err == unwrapError }

func Ok[T any](val T) Result[T] {
	return Result[T]{
		result: val,
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		err: err,
	}
}

func Errf[T any](format string, args ...any) Result[T] {
	return Result[T]{
		err: fmt.Errorf(format, args...),
	}
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return !r.IsOk()
}

func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panic(&ResultUnwrapError[T]{result: r.result, err: r.err})
	}
	return r.result
}

func (r Result[T]) UnwrapOrElse(back func(error) T) T {
	if r.IsErr() {
		return back(r.err)
	}
	return r.result
}

func (r Result[T]) UnwrapOrDefault(_default T) T {
	if r.IsErr() {
		return _default
	}
	return r.result
}

func (r Result[T]) UnwrapOrZeroValue() T {
	if r.IsErr() {
		t := new(T)
		return *t
	}
	return r.result
}

func (r Result[T]) Unzip() (T, error) {
	return r.result, r.err
}

func (r Result[T]) Error() error {
	return r.err
}
