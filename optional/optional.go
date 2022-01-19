package optional

import (
	"fmt"
)

type Optional[T any] struct {
	val *T
}

func Of[T any](sth T, ok bool) Optional[T] {
	if ok {
		return Some(sth)
	}
	return None[T]()
}

func Some[T any](sth T) Optional[T] {
	return Optional[T]{
		val: &sth,
	}
}

func None[T any]() Optional[T] {
	return Optional[T]{}
}

func (o Optional[T]) Ok() bool {
	return o.val != nil
}

func (o Optional[T]) Must() T {
	if !o.Ok() {
		panic(fmt.Errorf("option must return val of type: %T", o.val))
	}
	return *o.val
}

func (o Optional[T]) Value() T {
	if o.Ok() {
		return *o.val
	}
	t := new(T)
	return *t
}

func (o Optional[T]) ValueOr(sth T) T {
	if o.Ok() {
		return *o.val
	}
	return sth
}

func (o Optional[T]) Get() (T, bool) {
	return *o.val, o.Ok()
}

func (o Optional[T]) Prt() *T {
	return o.val
}
