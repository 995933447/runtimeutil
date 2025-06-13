package runtimeutil

import (
	"errors"
	"testing"
)

func TestStackErr(t *testing.T) {
	t.Log(testStackErr())
}

func testStackErr() error {
	return NewStackErr(errors.New("test error"))
}
