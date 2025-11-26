package runtimeutil

import (
	"bytes"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/petermattis/goid"
)

var goIdBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
}

var GetGoId = goid.Get

func GetGoIdFromStack() int64 {
	bp := goIdBuf.Get().(*[]byte)
	defer goIdBuf.Put(bp)
	buf := *bp
	return ExtractGID(buf[:runtime.Stack(buf[:], false)])
}

var goroutineSpaceOffset = len("goroutine ")

func ExtractGID(s []byte) int64 {
	s = s[goroutineSpaceOffset:]
	s = s[:bytes.IndexByte(s, ' ')]
	gid, _ := strconv.ParseInt(string(s), 10, 64)
	return gid
}

func IsGoVersionAtLeast(targetMajor, targetMinor, targetPatch int) bool {
	version := runtime.Version()
	if !strings.HasPrefix(version, "go") {
		return false
	}

	// 提取主版本、次版本和补丁版本
	re := regexp.MustCompile(`go(\d+)\.(\d+)(?:\.(\d+))?`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 3 {
		return false
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch := 0
	if len(matches) >= 4 && matches[3] != "" {
		patch, _ = strconv.Atoi(matches[3])
	}

	if major > targetMajor {
		return true
	} else if major == targetMajor {
		if minor > targetMinor {
			return true
		} else if minor == targetMinor {
			return patch >= targetPatch
		}
	}
	return false
}
