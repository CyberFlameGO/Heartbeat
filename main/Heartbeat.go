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

type Stats struct {
	User, System, Idle, Total uint64
	CpuCount, StatCount       uint64
}

type cpuStat struct {
	name string
	ptr  *uint64
}

func GetCpu() (*Stats, error) {
	file, err := os.Open("/proc/stat")

	if err != nil {
		fmt.Println("Could not read /proc/stat, found error: ", err.Error())
		return nil, err
	}

	defer file.Close()

	return getCpuStats(file)
}

func getCpuStats(out io.Reader) (*Stats, error) {
	scanner := bufio.NewScanner(out)
	var cpu Stats

	//assign all the data from the scanner to the stats in the struct

	cpuStats := []cpuStat{
		{"user", &cpu.User},
		{"system", &cpu.System},
		{"idle", &cpu.Idle},
		{"total", &cpu.Total},
	}

	if !scanner.Scan() {
		return nil, fmt.Errorf("Could not scan /proc/stat")
	}

	valStrs := strings.Fields(scanner.Text())[1:]

	cpu.StatCount = uint64(len(valStrs))

	for i, valStr := range valStrs {
		val, err := strconv.ParseUint(valStr, 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to scan section %s from /proc/stat", cpuStats[i].name)
		}

		*cpuStats[i].ptr = val
		cpu.Total += val

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
				cpu.CpuCount++
			}
		}

	}
	err := scanner.Err()

	if err != nil {
		return nil, fmt.Errorf("Scan error within /proc/stat %s", err)
	}

	return &cpu, nil
}

func main() {

}
