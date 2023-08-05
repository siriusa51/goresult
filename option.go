package goresult

type Option[T any] interface {
	Value() T
	IsSome() bool
	IsNone() bool
	Unwrap() T
	UnwrapOr(defaults T) T
	UnwrapOrElse(f func() T) T
	UnwrapOrDefault() T
	Inspect(f func(T)) Option[T]
	OkOr(err interface{}) Result[T]
	OkOrElse(f func() error) Result[T]
	Filter(predicate func(value T) bool) Option[T]
}

// option is an option type, it is either Some(T) or None.
type option[T any] struct {
	value T
	none  bool
}

// Some returns an option value of Some(T).
func Some[T any](value T) Option[T] {
	return &option[T]{value: value, none: false}
}

// None returns an option value of None.
func None[T any]() Option[T] {
	return &option[T]{none: true}
}

// Value return value
// example:
//
//	opt := Some(1)
//	fmt.Println(opt.Value())
func (opt *option[T]) Value() T {
	return opt.value
}

// IsSome returns true if the option is a Some value.
func (opt *option[T]) IsSome() bool {
	return !opt.none
}

// IsNone returns true if the option is a nil value.
func (opt *option[T]) IsNone() bool {
	return opt.none
}

// Unwrap returns the inner T of a Some(T). Panics if the self value equals nil.
// example:
//
//	opt := Some(1)
//	fmt.Println(opt.Unwrap())
//	// Output: 1
//
// opt := None[int]()
// fmt.Println(opt.Unwrap())
// // panic
func (opt *option[T]) Unwrap() T {
	if opt.IsNone() {
		panic("called `option.Unwrap()` on a `nil` value")
	}

	return opt.value
}

// UnwrapOr returns the inner T of a Some(T). Returns defaults if the self value is nil.
// example:
//
//	opt := Some(1)
//	fmt.Println(opt.UnwrapOr(2))
//	// Output: 1
//
// opt := None[int]()
// fmt.Println(opt.UnwrapOr(2))
// // Output: 2
func (opt *option[T]) UnwrapOr(defaults T) T {
	if opt.IsNone() {
		return defaults
	}

	return opt.value
}

// UnwrapOrElse calls f if the self value is nil. Returns the inner T of a Some(T).
// example:
//
//	opt := Some(1)
//	fmt.Println(opt.UnwrapOrElse(func() int {
//		return 2
//	}))
//	// Output: 1
//
// opt := None[int]()
//
//	fmt.Println(opt.UnwrapOrElse(func() int {
//		return 2
//	}))
//
// // Output: 2
func (opt *option[T]) UnwrapOrElse(f func() T) T {
	if opt.IsNone() {
		return f()
	}

	return opt.value
}

// UnwrapOrDefault returns the inner T of a Some(T). Returns the default value of T if the self value is nil.
// example:
//
//	opt := Some(1)
//	fmt.Println(opt.UnwrapOrDefault())
//	// Output: 1
//
// opt := None[int]()
// fmt.Println(opt.UnwrapOrDefault())
// // Output: 0
func (opt *option[T]) UnwrapOrDefault() T {
	if opt.IsNone() {
		return None[T]().Value()
	}

	return opt.value
}

// Inspect calls f if the self value equals Some(T).
// example:
//
//	option.Some(1).Inspect(func(v int) {
//		fmt.Println(v)
//	})
//	// Output: 1
//
//	option.None[int]().Inspect(func(v int) {
//		fmt.Println(v)
//	})
//
// // Output:
func (opt *option[T]) Inspect(f func(T)) Option[T] {
	if opt.IsSome() {
		f(opt.value)
	}

	return opt
}

// OkOr returns an Ok(T) containing the inner T of a Some(T).
// If the self value is nil, returns an Error(err) containing err.
// example:
//
//	opt := Some(1)
//	fmt.Println(opt.OkOr(errors.New("error")))
//	// Output: Ok(1)
//
// opt := None[int]()
// fmt.Println(opt.OkOr(errors.New("error")))
// // Output: Error(error)
func (opt *option[T]) OkOr(err interface{}) Result[T] {
	if opt.IsSome() {
		return Ok[T](opt.value)
	}

	return Error[T](err)
}

// OkOrElse returns an Ok(T) containing the inner T of a Some(T).
// If the self value is nil, calls f and returns an Error(err) containing the result.
// example:
//
//	opt := Some(1)
//	fmt.Println(opt.OkOrElse(func() error {
//		return errors.New("error")
//	}))
//	// Output: Ok(1)
//
// opt := None[int]()
//
//	fmt.Println(opt.OkOrElse(func() error {
//		return errors.New("error")
//	}))
//
// // Output: Error(error)
func (opt *option[T]) OkOrElse(f func() error) Result[T] {
	if opt.IsSome() {
		return Ok[T](opt.value)
	}

	return Error[T](f())
}

// Filter returns None if the self value equals None, otherwise calls predicate with the wrapped value and returns:
//   - Some(T) if predicate returns true (where T is the wrapped value)
//   - None() otherwise.
func (opt *option[T]) Filter(predicate func(value T) bool) Option[T] {
	if opt.IsSome() && predicate(opt.value) {
		return opt
	}

	return None[T]()
}
