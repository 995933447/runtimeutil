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

func NewStackErr(err error) *StackErr {
	stackErr := &StackErr{
		Err: err,
	}

	pc, callFile, callLine, ok := runtime.Caller(1)

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

func (s *StackErr) Error() string {
	return fmt.Sprintf("%s:%s:%d %s", s.FileName, s.CallFuncName, s.CallLine, s.Err.Error())
}
