package goresult

import (
	"fmt"
	"reflect"
)

// Result is a generic type that represents either success (Ok) or failure (Error).
type Result[T any] struct {
	Value T
	Error error
}

// Ok returns a Result that is Ok.
// example:
//
//	Ok(1)
//	Ok("hello")
//	Ok([]byte("world"))
//
// Ok[int32](10)
func Ok[T any](value T) *Result[T] {
	return &Result[T]{
		Value: value,
		Error: nil,
	}
}

// Error returns a Result that is Error.
// example:
//
//	Error[int](errors.New("something went wrong"))
//	Error[string]("something went wrong")
//	Error[any](fmt.Errorf("something went wrong"))
func Error[T any](err interface{}) *Result[T] {
	return &Result[T]{
		Error: covertError(err),
	}
}

// IsOk returns true if the result is Ok.
// example:
//
//	result := Ok(1)
//	fmt.Println(result.IsOk())
//
// // Output: true
func (r *Result[T]) IsOk() bool {
	return r.Error == nil
}

// IsError returns true if the result is Error.
// example:
//
//	result := Error[int](errors.New("something went wrong"))
//	fmt.Println(result.IsError())
//
// // Output: true
func (r *Result[T]) IsError() bool {
	return r.Error != nil
}

// Except returns the value if the result is Ok, otherwise it panics with the given message.
// example:
//
//	result := Ok(1)
//	fmt.Println(result.Except("something went wrong"))
//
// // Output: 1
//
//	result := Error[int](errors.New("something went wrong"))
//	fmt.Println(result.Except("something went wrong"))
//
// // panic: something went wrong
func (r *Result[T]) Except(msg string) T {
	if r.IsError() {
		unwrapErrorFailed(msg, r.Error)
	}

	return r.Value
}

// ExceptError returns the error if the result is Error, otherwise it panics with the given message.
// example:
//
//	result := Error[int](errors.New("something went wrong"))
//	fmt.Println(result.ExceptError("something went wrong"))
//
// // Output: something went wrong
//
//	result := Ok(1)
//	fmt.Println(result.ExceptError("something went wrong"))
//
// // panic: something went wrong
func (r *Result[T]) ExceptError(msg string) error {
	if r.IsOk() {
		unwrapValueFailed(msg, r.Value)
	}

	return r.Error
}

// Inspect calls the given function with the value if the result is Ok.
// example:
//
//	result := Ok(1)
//	result.Inspect(func(v int) {
//		fmt.Println(v)
//	})
//
// // Output: 1
//
//	result := Error[int](errors.New("something went wrong"))
//	result.Inspect(func(v int) {
//		fmt.Println(v)
//	})
//
// // Output:
func (r *Result[T]) Inspect(f func(T)) *Result[T] {
	if r.IsOk() {
		f(r.Value)
	}

	return r
}

// InspectError calls the given function with the error if the result is Error.
// example:
//
//	result := Error[int](errors.New("something went wrong"))
//	result.InspectError(func(err error) {
//		fmt.Println(err)
//	})
//
// // Output: something went wrong
//
//	result := Ok(1)
//	result.InspectError(func(err error) {
//		fmt.Println(err)
//	})
//
// // Output:
func (r *Result[T]) InspectError(f func(error)) *Result[T] {
	if r.IsError() {
		f(r.Error)
	}

	return r
}

// Unwrap returns the value if the result is Ok, otherwise it panics.
func (r *Result[T]) Unwrap() T {
	return r.Except("called `Result.Unwrap()` on an `Error` value")
}

// UnwrapError returns the error if the result is Error, otherwise it panics.
func (r *Result[T]) UnwrapError() error {
	return r.ExceptError("called `Result.UnwrapError()` on an `Value` value")
}

// UnwrapOr returns the value if the result is Ok, otherwise it returns the given default.
// example:
//
//	result := Ok(1)
//	fmt.Println(result.UnwrapOr(2))
//
// // Output: 1
//
//	result := Error[int](errors.New("something went wrong"))
//	fmt.Println(result.UnwrapOr(2))
//
// // Output: 2
func (r *Result[T]) UnwrapOr(defaults T) T {
	if r.IsOk() {
		return r.Value
	}

	return defaults
}

// UnwrapOrDefault returns the value if the result is Ok, otherwise it returns the default value of the type.
// example:
//
//	result := Ok(1)
//	fmt.Println(result.UnwrapOrDefault())
//
// // Output: 1
//
//	result := Error[int](errors.New("something went wrong"))
//	fmt.Println(result.UnwrapOrDefault())
//
// // Output: 0
func (r *Result[T]) UnwrapOrDefault() T {
	if r.IsOk() {
		return r.Value
	}

	return Result[T]{}.Value
}

// UnwrapOrElse returns the value if the result is Ok, otherwise it calls and returns the given function.
// example:
//
//	result := Ok(1)
//	fmt.Println(result.UnwrapOrElse(func() int {
//		return 2
//	}))
//
// // Output: 1
//
//	result := Error[int](errors.New("something went wrong"))
//	fmt.Println(result.UnwrapOrElse(func() int {
//		return 2
//	}))
//
// // Output: 2
func (r *Result[T]) UnwrapOrElse(f func() T) T {
	if r.IsOk() {
		return r.Value
	}

	return f()
}

// Option returns the value as an Option.
// - If the result is Ok, the returned Option will be Some(T) with the value.
// - If the result is Error, the returned Option will be None().
func (r *Result[T]) Option() *Option[T] {
	if r.IsOk() {
		return Some(r.Value)
	}

	return None[T]()
}

func unwrapErrorFailed[E error](msg string, err E) {
	panic(fmt.Errorf("%s: %s", msg, err.Error()))
}

func unwrapValueFailed[T any](msg string, value T) {
	types := reflect.TypeOf(value).String()
	panic(fmt.Errorf("%s: %s", msg, types))
}
