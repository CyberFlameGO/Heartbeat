package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DiskStats struct {
	Name                            string
	ReadsCompleted, WritesCompleted uint64
}

func GetDiskStats() ([]DiskStats, error) {
	file, err := os.Open("/proc/diskstats")

	if err != nil {
		return nil, fmt.Errorf("Could not read /proc/diskstats, found error: %s", err)
	}

	defer file.Close()
	return getDiskStats(file)
}

func getDiskStats(out io.Reader) ([]DiskStats, error) {
	scanner := bufio.NewScanner(out)
	var diskStats []DiskStats

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) < 14 {
			continue
		}

		name := fields[2]
		readsCompleted, err := strconv.ParseUint(fields[3], 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse reads completed, found error: %s", err)
		}

		writesCompleted, err := strconv.ParseUint(fields[7], 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse writes completed, found error: %s", err)
		}

		diskStats = append(diskStats, DiskStats{
			Name:            name,
			WritesCompleted: writesCompleted,
			ReadsCompleted:  readsCompleted,
		})
	}

	err := scanner.Err()

	if err != nil {
		return nil, fmt.Errorf("Could not use the scanner, found error: %s", err)
	}

	return diskStats, nil
}
