package optional

import (
	"testing"

	"github.com/nickchenyx/ggwp/internal/assert"
)

func TestOptionalValue(t *testing.T) {
	assert.Equal(t, 10, Some(10).Value())
	assert.Equal(t, 10, Some(10).Value())
	assert.Equal(t, 0, None[int]().Value())
	assert.Equal(t, 10, Of(10, true).Value())
	assert.Equal(t, 0, Of(10, false).Value()) // NOTE: not recommend
	assert.Equal(t, 0, Of(0, false).Value())
}

func TestOptionalValueOr(t *testing.T) {
	assert.Equal(t, 10, Some(10).ValueOr(1))
	assert.Equal(t, 1, None[int]().ValueOr(1))
	assert.Equal(t, 10, Of(10, true).ValueOr(1))
	assert.Equal(t, 1, Of(10, false).ValueOr(1)) // NOTE: not recommend
	assert.Equal(t, 1, Of(0, false).ValueOr(1))
}

func TestOptionalValueOk(t *testing.T) {
	assert.True(t, Some(10).Ok())
	assert.True(t, Some(0).Ok())
	assert.False(t, None[int]().Ok())
	assert.True(t, Of(10, true).Ok())
	assert.False(t, Of(10, false).Ok()) // NOTE: not recommend
	assert.False(t, Of(0, false).Ok())
}

func TestOptionalValueMust(t *testing.T) {
	assert.NotPanic(t, func() { Some(10).Must() })
	assert.NotPanic(t, func() { Some(0).Must() })
	assert.Panic(t, func() { None[int]().Must() })
	assert.NotPanic(t, func() { Of(10, true).Must() })
	assert.Panic(t, func() { Of(10, false).Must() }) // NOTE: not recommend
	assert.Panic(t, func() { Of(0, false).Must() })
}
