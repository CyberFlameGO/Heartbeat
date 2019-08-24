package main

import (
	"fmt"
	"io"
	"os"
)

type MemoryStats struct {
	Total, Used, Buffers, Cached, Free, Available, Active, Inactive uint64
}

func GetMemory() (*MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")

	if err != nil {
		return nil, fmt.Errorf("Could not open /proc/meminfo, found error: %s", err)
	}

	defer file.Close()
	return getMemoryStats(file)
}

func getMemoryStats(out io.Reader) (*MemoryStats, error) {

}
