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

type StatsCpu struct {
	User, System, Idle, Total uint64
	CpuCount, StatCount       uint64
}

type cpuStat struct {
	name string
	ptr  *uint64
}

func GetCpu() (*StatsCpu, error) {
	file, err := os.Open("/proc/stat")

	if err != nil {
		fmt.Println("Could not read /proc/stat, found error: ", err.Error())
		return nil, err
	}

	defer file.Close()

	return getCpuStats(file)
}

func getCpuStats(out io.Reader) (*StatsCpu, error) {
	scanner := bufio.NewScanner(out)
	var cpu StatsCpu

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

type StatsDisk struct {
	Name                            string
	ReadsCompleted, WritesCompleted uint64
}

func GetDisk() ([]StatsDisk, error) {
	file, err := os.Open("/proc/diskstats")

	if err != nil {
		fmt.Println("Could not read /proc/diskstats, found error: ", err)
		return nil, err
	}

	defer file.Close()

	return getDiskStats(file)
}

func getDiskStats(out io.Reader) ([]StatsDisk, error) {
	scanner := bufio.NewScanner(out)
	var diskStats []StatsDisk

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) < 14 {
			continue
		}

		name := fields[2]
		readsCompleted, err := strconv.ParseUint(fields[3], 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse reads completed of %s", name)
		}

		writesCompleted, err := strconv.ParseUint(fields[7], 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse writes completed of %s", name)
		}

		diskStats = append(diskStats, StatsDisk{
			Name:            name,
			WritesCompleted: writesCompleted,
			ReadsCompleted:  readsCompleted,
		})

	}

	err := scanner.Err()

	if err != nil {
		return nil, fmt.Errorf("Scan error for /proc/diskstats %s", err)
	}

	return diskStats, nil
}

func main() {

}
