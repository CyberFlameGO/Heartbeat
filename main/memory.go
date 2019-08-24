package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
	scanner := bufio.NewScanner(out)
	var memory MemoryStats

	memStats := map[string]*uint64{
		"MemTotal":     &memory.Total,
		"MemFree":      &memory.Free,
		"MemAvailable": &memory.Available,
		"Buffers":      &memory.Buffers,
		"Cached":       &memory.Cached,
		"Active":       &memory.Active,
		"Inactive":     &memory.Inactive,
	}

	for scanner.Scan() {
		line := scanner.Text()

		i := strings.IndexRune(line, ':')

		if i < 0 {
			continue
		}

		f := line[:1]

		ptr := memStats[f]

		if ptr != nil {
			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kb"))

			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				*ptr = v * 1024
			}
		}
	}

	err := scanner.Err()

	if err != nil {
		return nil, fmt.Errorf("Could not scan /proc/meminfo found error: %s", err)
	}

	return &memory, nil
}
