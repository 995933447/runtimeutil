package runtimeutil

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func GetMemoryByRuntime() uint64 {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	return stat.Alloc
}

func GetCurrMemory() uint64 {
	mem, _ := GetMemoryByProcStatusFile()
	if mem > 0 {
		return mem
	}

	return GetMemoryByRuntime()
}

func GetMemoryByProcStatusFile() (uint64, error) {
	data, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memKB, err := strconv.ParseUint(fields[1], 10, 64) // 单位为 KB
				if err != nil {
					return 0, err
				}

				return memKB * 1024, nil
			}
		}
	}

	return 0, fmt.Errorf("VmRSS not found")
}
