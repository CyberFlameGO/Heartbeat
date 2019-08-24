package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type CPUStats struct {
	User, System, Idle, Total, CPUCount uint64
}

type CPUInfo struct {
	name string
	ptr  *uint64
}

func GetCpuStats() (*CPUStats, error) {
	file, err := os.Open("/proc/stat")

	if err != nil {
		_ = fmt.Errorf("Could not read /proc/stat, found error: %s", err)
		return nil, err
	}

	defer file.Close()

	return getCpuStats(file)
}

func getCpuStats(out io.Reader) (*CPUStats, error) {
	scanner := bufio.NewScanner(out)

	var cpu CPUStats

	cpuStats := []CPUInfo{
		{"user", &cpu.User},
		{"system", &cpu.System},
		{"idle", &cpu.Idle},
		{"total", &cpu.Total},
	}

	if !scanner.Scan() {
		return nil, fmt.Errorf("Could not scan /proc/stat")
	}

	valStrs := strings.Fields(scanner.Text())[1:]

	for i, valStr := range valStrs {
		val, err := strconv.ParseUint(valStr, 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to scan section %s ", cpuStats[i].name)
		}

		*cpuStats[i].ptr = val
		cpu.Total += val

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
				cpu.CPUCount++
			}
		}
	}

	err := scanner.Err()

	if err != nil {
		return nil, fmt.Errorf("Scan error with /proc/stat {%s}", err)
	}

	return &cpu, nil
}
