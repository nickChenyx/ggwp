package result

import (
	"fmt"
	"os"
	"testing"

	"github.com/nickchenyx/ggwp/internal/assert"
)

func TestCallWithWrapPanic(t *testing.T) {
	assert.True(t, CallWithWrapPanic(func() int {
		panic(fmt.Errorf("panic err"))
	}).IsErr())
	assert.NotPanic(t, func() {
		CallWithWrapPanic(func() int {
			panic("panic here")
		})
	})
	assert.True(t, CallWithWrapPanic(func() int {
		panic("panic here")
	}).IsErr())
}

func TestCall(t *testing.T) {
	assert.Panic(t, func() {
		Call(func() int {
			panic("panic here")
		})
	})

	assert.NotPanic(t, func() {
		Call(func() int {
			return Err[int](fmt.Errorf("none")).Unwrap()
		})
	})

	assert.True(t, Call(func() int {
		return Err[int](fmt.Errorf("none")).Unwrap()
	}).IsErr())

	err := os.ErrClosed
	resultError := Err[int](err)

	assert.Equal(t, err, Call(func() int {
		return resultError.Unwrap()
	}).err)
}

func TestCallExportPanicCustom(t *testing.T) {
	err := os.ErrClosed
	assert.True(t, CallExportPanicCustom(func() int {
		panic(err)
	}, func(e error) bool {
		return e == err
	}).IsErr())

	resultError := Err[int](err)

	assert.Equal(t, err, CallExportPanicCustom(func() int {
		return resultError.Unwrap()
	}, nil).err)

	assert.Panic(t, func() {
		CallExportPanicCustom(func() int {
			panic(fmt.Errorf("panic err"))
		}, func(e error) bool {
			return false
		})
	})

	assert.Panic(t, func() {
		CallExportPanicCustom(func() int {
			panic("panic here")
		}, func(e error) bool {
			return false
		})
	})
}
