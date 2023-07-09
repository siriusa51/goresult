package goresult

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertError(t *testing.T) {
	assert.Error(t, covertError("error"))
	assert.Error(t, covertError(1))
	assert.Error(t, covertError(fmt.Errorf("error")))
}
