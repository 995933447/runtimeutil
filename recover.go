package runtimeutil

import (
	"fmt"
	"gorm.io/gorm/logger"
	"runtime/debug"
	"strings"
	"sync"
)

var (
	Recover = RecoverToPrint
	accessRecoverMu sync.RWMutex
)

func SetRecover(recover func()) {
	accessRecoverMu.Lock()
	defer accessRecoverMu.Unlock()
	Recover = recover
}

func RecoverToPrint() {
	if err := recover(); err != nil {
		stack := debug.Stack()
		if len(stack) > 0 {
			lines := strings.Split(string(stack), "\n")
			for _, line := range lines {
				fmt.Println(line)
			}
			return
		}
		fmt.Printf("stack is empty (%s)\n", err)
	}
}
