package runtimeutil

import (
	"errors"
	"testing"
)

func TestStackErr(t *testing.T) {
	t.Log(testStackErr())
	t.Log(testStackErrWithSkip(1))
	t.Log(testStackErrWithSkip(2))
}

func testStackErr() error {
	return NewStackErr(errors.New("test error"))
}

func testStackErrWithSkip(skip int) error {
	return NewStackErrWithSkip(skip, errors.New("test error"))
}
