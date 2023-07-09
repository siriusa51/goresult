package goresult

import "fmt"

func covertError(err interface{}) error {
	switch err.(type) {
	case error:
		return err.(error)
	default:
		return fmt.Errorf("%v", err)
	}
}
