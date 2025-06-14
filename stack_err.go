package runtimeutil

import (
	"fmt"
	"path"
	"runtime"
)

type StackErr struct {
	Err          error
	CallFuncName string
	FileName     string
	CallLine     int
}

func NewStackErrWithSkip(skip int, err error) *StackErr {
	stackErr := &StackErr{
		Err: err,
	}

	pc, callFile, callLine, ok := runtime.Caller(skip)

	var callFuncName string
	if ok {
		callFuncName = runtime.FuncForPC(pc).Name()
	}

	_, callFile = path.Split(callFile)

	stackErr.CallFuncName = callFuncName
	stackErr.FileName = callFile
	stackErr.CallLine = callLine

	return stackErr
}

func NewStackErr(err error) *StackErr {
	return NewStackErrWithSkip(2, err)
}

func (s *StackErr) Error() string {
	return fmt.Sprintf("%s:%s:%d line: %s", s.FileName, s.CallFuncName, s.CallLine, s.Err.Error())
}
