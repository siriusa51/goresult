package goresult

import (
	"fmt"
	"reflect"
)

type Result[T any] interface {
	Value() T
	Error() error
	IsOk() bool
	IsError() bool
	Except(msg string) T
	ExceptError(msg string) error
	Inspect(f func(T)) Result[T]
	InspectError(f func(error)) Result[T]
	Unwrap() T
	UnwrapError() error
	UnwrapOr(defaults T) T
	UnwrapOrDefault() T
	UnwrapOrElse(f func() T) T
	Option() Option[T]
}

// result is a generic type that represents either success (Ok) or failure (Error).
type result[T any] struct {
	value T
	error error
}

// Ok returns a result that is Ok.
// example:
//
//	Ok(1)
//	Ok("hello")
//	Ok([]byte("world"))
//
// Ok[int32](10)
func Ok[T any](value T) Result[T] {
	return &result[T]{
		value: value,
		error: nil,
	}
}

// Error returns a result that is Error.
// example:
//
//	error[int](errors.New("something went wrong"))
//	error[string]("something went wrong")
//	error[any](fmt.Errorf("something went wrong"))
func Error[T any](err interface{}) Result[T] {
	return &result[T]{
		error: covertError(err),
	}
}

// Value return value
// example:
//
//	result := Ok(1)
//	fmt.Println(result.value())
//
// // Output: 1
func (r *result[T]) Value() T {
	return r.value
}

// Error return error
// example:
//
//	result := error[int](errors.New("something went wrong"))
//	fmt.Println(result.Error())
//
// // Output: something went wrong
func (r *result[T]) Error() error {
	return r.error
}

// IsOk returns true if the result is Ok.
// example:
//
//	result := Ok(1)
//	fmt.Println(result.IsOk())
//
// // Output: true
func (r *result[T]) IsOk() bool {
	return r.error == nil
}

// IsError returns true if the result is Error.
// example:
//
//	result := error[int](errors.New("something went wrong"))
//	fmt.Println(result.IsError())
//
// // Output: true
func (r *result[T]) IsError() bool {
	return r.error != nil
}

// Except returns the value if the result is Ok, otherwise it panics with the given message.
// example:
//
//	result := Ok(1)
//	fmt.Println(result.Except("something went wrong"))
//
// // Output: 1
//
//	result := error[int](errors.New("something went wrong"))
//	fmt.Println(result.Except("something went wrong"))
//
// // panic: something went wrong
func (r *result[T]) Except(msg string) T {
	if r.IsError() {
		unwrapErrorFailed(msg, r.error)
	}

	return r.value
}

// ExceptError returns the error if the result is Error, otherwise it panics with the given message.
// example:
//
//	result := error[int](errors.New("something went wrong"))
//	fmt.Println(result.ExceptError("something went wrong"))
//
// // Output: something went wrong
//
//	result := Ok(1)
//	fmt.Println(result.ExceptError("something went wrong"))
//
// // panic: something went wrong
func (r *result[T]) ExceptError(msg string) error {
	if r.IsOk() {
		unwrapValueFailed(msg, r.value)
	}

	return r.error
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
//	result := error[int](errors.New("something went wrong"))
//	result.Inspect(func(v int) {
//		fmt.Println(v)
//	})
//
// // Output:
func (r *result[T]) Inspect(f func(T)) Result[T] {
	if r.IsOk() {
		f(r.value)
	}

	return r
}

// InspectError calls the given function with the error if the result is Error.
// example:
//
//	result := error[int](errors.New("something went wrong"))
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
func (r *result[T]) InspectError(f func(error)) Result[T] {
	if r.IsError() {
		f(r.error)
	}

	return r
}

// Unwrap returns the value if the result is Ok, otherwise it panics.
func (r *result[T]) Unwrap() T {
	return r.Except("called `result.Unwrap()` on an `error` value")
}

// UnwrapError returns the error if the result is Error, otherwise it panics.
func (r *result[T]) UnwrapError() error {
	return r.ExceptError("called `result.UnwrapError()` on an `value` value")
}

// UnwrapOr returns the value if the result is Ok, otherwise it returns the given default.
// example:
//
//	result := Ok(1)
//	fmt.Println(result.UnwrapOr(2))
//
// // Output: 1
//
//	result := error[int](errors.New("something went wrong"))
//	fmt.Println(result.UnwrapOr(2))
//
// // Output: 2
func (r *result[T]) UnwrapOr(defaults T) T {
	if r.IsOk() {
		return r.value
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
//	result := error[int](errors.New("something went wrong"))
//	fmt.Println(result.UnwrapOrDefault())
//
// // Output: 0
func (r *result[T]) UnwrapOrDefault() T {
	if r.IsOk() {
		return r.value
	}

	return result[T]{}.value
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
//	result := error[int](errors.New("something went wrong"))
//	fmt.Println(result.UnwrapOrElse(func() int {
//		return 2
//	}))
//
// // Output: 2
func (r *result[T]) UnwrapOrElse(f func() T) T {
	if r.IsOk() {
		return r.value
	}

	return f()
}

// option returns the value as an option.
// - If the result is Ok, the returned option will be Some(T) with the value.
// - If the result is Error, the returned option will be None().
func (r *result[T]) Option() Option[T] {
	if r.IsOk() {
		return Some(r.value)
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
