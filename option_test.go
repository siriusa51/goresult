package goresult

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Option_Some(t *testing.T) {
	opt := Some(1)

	assert.Equal(t, opt.IsSome(), true)
	assert.Equal(t, opt.IsNone(), false)
	assert.Equal(t, opt.Value(), 1)
}

func Test_Option_None(t *testing.T) {
	opt := None[int]()

	assert.Equal(t, opt.IsSome(), false)
	assert.Equal(t, opt.IsNone(), true)
	assert.Equal(t, opt.Value(), 0)
}

func Test_Option_Unwrap(t *testing.T) {
	opt := Some(1)

	assert.Equal(t, opt.Unwrap(), 1)
}

func Test_Option_Unwrap_None(t *testing.T) {
	opt := None[int]()

	assert.Panics(t, func() { opt.Unwrap() }, "called `option.Unwrap()` on a `nil` value")
}

func Test_Option_UnwrapOr(t *testing.T) {
	opt := Some(1)

	assert.Equal(t, opt.UnwrapOr(2), 1)
}

func Test_Option_UnwrapOr_None(t *testing.T) {
	opt := None[int]()

	assert.Equal(t, opt.UnwrapOr(2), 2)
}

func Test_Option_UnwrapOrElse(t *testing.T) {
	opt := Some(1)

	assert.Equal(t, opt.UnwrapOrElse(func() int { return 2 }), 1)
}

func Test_Option_UnwrapOrElse_None(t *testing.T) {
	opt := None[int]()

	assert.Equal(t, opt.UnwrapOrElse(func() int { return 2 }), 2)
}

func Test_Option_UnwrapOrDefault(t *testing.T) {
	opt := Some(1)

	assert.Equal(t, opt.UnwrapOrDefault(), 1)
}

func Test_Option_UnwrapOrDefault_None(t *testing.T) {
	opt := None[int]()

	assert.Equal(t, opt.UnwrapOrDefault(), 0)
}

func Test_Option_Inspect(t *testing.T) {
	opt := Some(1)

	var result int
	opt.Inspect(func(i int) { result = i })

	assert.Equal(t, result, 1)
}

func Test_Option_Inspect_None(t *testing.T) {
	opt := None[int]()

	var result int
	opt.Inspect(func(i int) { result = i })

	assert.Equal(t, result, 0)
}

func Test_Option_InspectError(t *testing.T) {
	opt := Some(1)

	var result int
	opt.Inspect(func(i int) { result = i })

	assert.Equal(t, result, 1)
}

func Test_Option_InspectError_None(t *testing.T) {
	opt := None[int]()

	var result int
	opt.Inspect(func(i int) { result = i })

	assert.Equal(t, result, 0)
}

func Test_Option_OkOr(t *testing.T) {
	opt := Some(1)

	assert.Equal(t, opt.OkOr(fmt.Errorf("")), Ok(1))
}

func Test_Option_OkOr_None(t *testing.T) {
	opt := None[string]()

	assert.Equal(t, opt.OkOr(fmt.Errorf("")), Error[string](fmt.Errorf("")))
}

func Test_Option_OkOrElse(t *testing.T) {
	opt := Some(1)

	assert.Equal(t, opt.OkOrElse(func() error { return fmt.Errorf("") }), Ok(1))
}

func Test_Option_OkOrElse_None(t *testing.T) {
	opt := None[string]()

	assert.Equal(t, opt.OkOrElse(func() error { return fmt.Errorf("") }), Error[string](fmt.Errorf("")))
}

func Test_Option_Filter(t *testing.T) {
	opt := Some(1)

	assert.Equal(t, opt.Filter(func(i int) bool { return i == 1 }), opt)
}

func Test_Option_Filter_None(t *testing.T) {
	opt := None[int]()

	assert.Equal(t, opt.Filter(func(i int) bool { return i == 1 }), opt)
}
