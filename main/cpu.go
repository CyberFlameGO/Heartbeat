package main

import (
	"fmt"
	"io"
	"os"
)

type CPUStats struct {
	User, System, Idle, Total uint64
	CPUCount, StatCount       uint64
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

}
