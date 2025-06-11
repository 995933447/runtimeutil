package runtimeutil

import (
	"encoding/binary"
	"github.com/petermattis/goid"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789xy"
)

var pid = uint16(0)

var TraceMap sync.Map

func GenTraceStrId() string {
	if pid == 0 {
		pid = uint16(os.Getpid())
	}
	var buf [15]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(time.Now().UnixNano()))
	binary.LittleEndian.PutUint16(buf[8:], pid)
	binary.LittleEndian.PutUint32(buf[10:], uint32(rand.Int()))
	//for i := 1; i < 7; i++ {
	//	x := buf[i]
	//	buf[i] = buf[14-i]
	//	buf[14-i] = x
	//}
	// 3 byte -> 4 char
	var s [20]byte
	d := 0
	for i := 0; i < 15; {
		val := uint(buf[i])<<16 | uint(buf[i+1])<<8 | uint(buf[i+2])
		i += 3
		s[d] = chars[val>>18&0x3F]
		s[d+1] = chars[val>>12&0x3F]
		s[d+2] = chars[val>>6&0x3F]
		s[d+3] = chars[val&0x3F]
		d += 4
	}
	return string(s[:])
}

func GetOrCreateTrace() string {
	gpid := goid.Get()
	traceAny, ok := TraceMap.Load(gpid)
	if ok {
		return traceAny.(string)
	}

	trace := GenTraceStrId()
	TraceMap.Store(gpid, trace)
	return trace
}

func StoreTrace(trace string) {
	TraceMap.Store(goid.Get(), trace)
}

func RemoveTrace(gid int64) {
	TraceMap.Delete(gid)
}

func AutoRemoveTrace() {
	RemoveTrace(goid.Get())
}
