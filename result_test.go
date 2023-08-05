package goresult

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

type testType struct {
	String string
	Int    int
	Point  *float64
	Array  []byte
	Map    map[string]interface{}
}

func newTestType() (testType, testType) {
	v1 := strconv.Itoa(rand.Int())
	v2 := rand.Int()
	v3 := rand.Float64()
	v4 := strconv.Itoa(rand.Int())
	v5 := map[string]interface{}{
		"v1": v1,
		"v2": v2,
		"v3": v3,
		"v4": v4,
	}

	v6 := map[string]interface{}{
		"v1": v1,
		"v2": v2,
		"v3": v3,
		"v4": v4,
	}

	return testType{v1, v2, &v3, []byte(v4), v5}, testType{v1, v2, &v3, []byte(v4), v6}
}

func equal[T any](t *testing.T, active, excepted T) {
	assert.Equal(t, excepted, active, "Expected value to be %v, but got %v", excepted, active)
}

func isNil[T any](t *testing.T, active T) {
	assert.Nil(t, active, "Expected value to be nil, but got %v", active)
}

func assertOk[T any](t *testing.T, result Result[T], excepted T) {
	typeName := reflect.TypeOf(excepted).String()
	t.Run(fmt.Sprintf("Ok(%s)", typeName), func(t *testing.T) {
		equal(t, result.Value(), excepted)
		isNil(t, result.Error())
	})
}

func Test_Result_Ok(t *testing.T) {
	assertOk(t, Ok(42), 42)
	assertOk(t, Ok("hello"), "hello")
	assertOk(t, Ok([]byte("world")), []byte("world"))
	assertOk(t, Ok(1.1), 1.1)
	assertOk(t, Ok(true), true)
	assertOk(t, Ok([]string{"v1", "v2"}), []string{"v1", "v2"})
	assertOk(t, Ok(map[string]interface{}{"value": 1}), map[string]interface{}{"value": 1})
	assertOk(t, Ok(fmt.Errorf("no error")), fmt.Errorf("no error"))

	v1, v2 := newTestType()

	assertOk(t, Ok(v1), v2)

	v1, v2 = newTestType()
	assertOk(t, Ok(&v2), &v2)
}

func assertError[T any](t *testing.T, result Result[T], excepted error) {
	typeName := reflect.TypeOf(result.Value()).String()
	t.Run(fmt.Sprintf("error(%s)", typeName), func(t *testing.T) {
		equal(t, result.Error(), excepted)
	})
}

func Test_Result_Error(t *testing.T) {
	assertError(t, Error[int](fmt.Errorf("error")), fmt.Errorf("error"))
	assertError(t, Error[int]("error int"), fmt.Errorf("error int"))
	assertError(t, Error[float64]("error float64"), fmt.Errorf("error float64"))
	assertError(t, Error[string]("error string"), fmt.Errorf("error string"))
	assertError(t, Error[[]byte]("error []byte"), fmt.Errorf("error []byte"))
	assertError(t, Error[map[string]interface{}]("error map[string]interface{}"), fmt.Errorf("error map[string]interface{}"))
}

func assertIsOk[T any](t *testing.T, result Result[T], isOk bool) {
	typeName := reflect.TypeOf(isOk).String()
	t.Run(fmt.Sprintf("IsOk(%s)", typeName), func(t *testing.T) {
		equal(t, result.IsOk(), isOk)
		if isOk {
			isNil(t, result.Error())
		}
	})
}

func Test_Result_IsOk(t *testing.T) {
	assertIsOk(t, Ok(1), true)
	assertIsOk(t, Ok("string"), true)
	assertIsOk(t, Ok([]byte("hello")), true)
	assertIsOk(t, Error[[]byte]([]byte("hello")), false)
	assertIsOk(t, Error[float64](fmt.Errorf("error")), false)
}

func assertIsError[T any](t *testing.T, result Result[T], isError bool) {
	typeName := reflect.TypeOf(isError).String()
	t.Run(fmt.Sprintf("IsError(%s)", typeName), func(t *testing.T) {
		equal(t, result.IsError(), isError)
		if !isError {
			isNil(t, result.Error())
		}
	})
}

func Test_Result_IsError(t *testing.T) {
	assertIsError(t, Ok(1), false)
	assertIsError(t, Ok("string"), false)
	assertIsError(t, Ok([]byte("hello")), false)
	assertIsError(t, Error[[]byte]([]byte("hello")), true)
	assertIsError(t, Error[float64](fmt.Errorf("error")), true)
}

func assertExcept[T any](t *testing.T, result Result[T], excepted T, isPanic bool) {
	typeName := reflect.TypeOf(excepted).String()
	t.Run(fmt.Sprintf("Except(%s)", typeName), func(t *testing.T) {
		if isPanic {
			assert.Panics(t, func() {
				val := result.Except("error")
				equal(t, val, excepted)
			}, "Expected panic, but not")
		} else {
			equal(t, result.Except("error"), excepted)
		}
	})
}

func Test_Result_Except(t *testing.T) {
	assertExcept(t, Ok(42), 42, false)
	assertExcept(t, Ok("hello"), "hello", false)
	assertExcept(t, Ok([]byte("world")), []byte("world"), false)

	assertExcept(t, Error[int]("error"), 0, true)
	assertExcept(t, Error[string]("error"), "", true)
	assertExcept(t, Error[[]byte]("error"), nil, true)
}

func assertExceptError[T any](t *testing.T, result Result[T], isPanic bool) {
	typeName := reflect.TypeOf(result.Value()).String()
	t.Run(fmt.Sprintf("ExceptError(%s)", typeName), func(t *testing.T) {
		if isPanic {
			assert.Panics(t, func() {
				_ = result.ExceptError("error")

			}, "Expected panic, but not")
		} else {
			err := result.ExceptError("error")
			assert.Error(t, err, "Expected error, but got %v", err)
		}
	})
}

func Test_Result_ExceptError(t *testing.T) {
	assertExceptError(t, Ok(42), true)
	assertExceptError(t, Ok("hello"), true)
	assertExceptError(t, Ok([]byte("world")), true)

	assertExceptError(t, Error[int]("error"), false)
	assertExceptError(t, Error[string]("error"), false)
	assertExceptError(t, Error[[]byte]("error"), false)
}

func assertInspect[T any](t *testing.T, result Result[T], ok bool) {
	typeName := reflect.TypeOf(result.Value()).String()
	t.Run(fmt.Sprintf("Inspect(%s)", typeName), func(t *testing.T) {
		count := 1
		result.Inspect(func(_ T) {
			count += 1
		})
		if ok {
			assert.Equal(t, 2, count, "Expected call func(){count+=1} but not")
		} else {
			assert.Equal(t, 1, count, "Expected count not change, but changed")
		}
	})
}

func Test_Result_Inspect(t *testing.T) {
	assertInspect(t, Ok(42), true)
	assertInspect(t, Ok("hello"), true)
	assertInspect(t, Ok([]byte("world")), true)

	assertInspect(t, Error[int]("error"), false)
}

func assertInspectError[T any](t *testing.T, result Result[T], err bool) {
	typeName := reflect.TypeOf(result.Value()).String()
	t.Run(fmt.Sprintf("InspectError(%s)", typeName), func(t *testing.T) {
		count := 1
		result.InspectError(func(_ error) {
			count += 1
		})
		if err {
			assert.Equal(t, 2, count, "Expected call func(){count+=1} but not")
		} else {
			assert.Equal(t, 1, count, "Expected count not change, but changed")
		}
	})
}

func Test_Result_InspectError(t *testing.T) {
	assertInspectError(t, Ok(42), false)
	assertInspectError(t, Ok("hello"), false)
	assertInspectError(t, Ok([]byte("world")), false)

	assertInspectError(t, Error[int]("error"), true)
}

func Test_Result_Unwrap(t *testing.T) {
	assert.Equal(t, 42, Ok(42).Unwrap())
	assert.Equal(t, "hello", Ok("hello").Unwrap())
	assert.Equal(t, []byte("world"), Ok([]byte("world")).Unwrap())

	assert.Panics(t, func() {
		Error[int]("error").Unwrap()
	}, "Expected panic, but not")
}

func Test_Result_UnwrapError(t *testing.T) {
	assert.Panics(t, func() {
		_ = Ok(42).UnwrapError()
	}, "Expected panic, but not")

	assert.Panics(t, func() {
		_ = Ok("hello").UnwrapError()
	}, "Expected panic, but not")
}

func Test_Result_UnwrapOr(t *testing.T) {
	assert.Equal(t, 42, Ok(42).UnwrapOr(0))
	assert.Equal(t, "hello", Ok("hello").UnwrapOr(""))
	assert.Equal(t, []byte("world"), Ok([]byte("world")).UnwrapOr(nil))

	assert.Equal(t, 0, Error[int]("error").UnwrapOr(0))
	assert.Equal(t, "", Error[string]("error").UnwrapOr(""))
	assert.Nil(t, Error[[]byte]("error").UnwrapOr(nil))
}

func Test_Result_UnwrapOrDefault(t *testing.T) {
	assert.Equal(t, 42, Ok(42).UnwrapOrDefault())
	assert.Equal(t, "hello", Ok("hello").UnwrapOrDefault())
	assert.Equal(t, []byte("world"), Ok([]byte("world")).UnwrapOrDefault())

	assert.Equal(t, 0, Error[int]("error").UnwrapOrDefault())
	assert.Equal(t, "", Error[string]("error").UnwrapOrDefault())
	assert.Nil(t, Error[[]byte]("error").UnwrapOrDefault())
}

func Test_Result_UnwrapOrElse(t *testing.T) {
	assert.Equal(t, 42, Ok(42).UnwrapOrElse(func() int { return 0 }))
	assert.Equal(t, "hello", Ok("hello").UnwrapOrElse(func() string { return "" }))
	assert.Equal(t, []byte("world"), Ok([]byte("world")).UnwrapOrElse(func() []byte { return nil }))

	assert.Equal(t, 0, Error[int]("error").UnwrapOrElse(func() int { return 0 }))
	assert.Equal(t, "", Error[string]("error").UnwrapOrElse(func() string { return "" }))
	assert.Nil(t, Error[[]byte]("error").UnwrapOrElse(func() []byte { return nil }))
}

func Test_Result_Option(t *testing.T) {
	assert.Equal(t, Some(42), Ok(42).Option())
	assert.Equal(t, Some("hello"), Ok("hello").Option())
	assert.Equal(t, Some([]byte("world")), Ok([]byte("world")).Option())

	assert.Equal(t, None[int](), Error[int]("error").Option())
	assert.Equal(t, None[string](), Error[string]("error").Option())
	assert.Equal(t, None[[]byte](), Error[[]byte]("error").Option())
}

func Test_unwrapErrorFailed(t *testing.T) {
	assert.Panics(t, func() {
		unwrapErrorFailed[error]("err", fmt.Errorf("error"))
	}, "Expected panic, but not")
}

func Test_unwrapValueFailed(t *testing.T) {
	assert.Panics(t, func() {
		unwrapValueFailed[string]("err", "value")
	}, "Expected panic, but not")
}
