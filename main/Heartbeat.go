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

func GetCpu() *StatsCpu {
	file, err := os.Open("/proc/stat")

	if err != nil {
		fmt.Println("Could not read /proc/stat, found error: ", err.Error())
		return nil
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

type StatsMemory struct {
	Total, Used, Buffers, Cached, Free, Available, Active, Inactive uint64
	MemAvailableEnabled                                             bool
}

func GetMemory() (*StatsMemory, error) {
	file, err := os.Open("/proc/meminfo")

	if err != nil {
		return nil, err
	}

	defer file.Close()
	return getMemoryStats(file)
}

func getMemoryStats(out io.Reader) (*StatsMemory, error) {
	scanner := bufio.NewScanner(out)
	var memory StatsMemory

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

		f := line[:i]

		if ptr := memStats[f]; ptr != nil {
			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kb"))

			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				*ptr = v * 1024
			}

			if f == "MemAvailable" {
				memory.MemAvailableEnabled = true
			}
		}
	}
	err := scanner.Err()

	if err != nil {
		return nil, fmt.Errorf("Scan error for /proc/meminfo %s", err)
	}

	return &memory, nil
}

func main() {

	var cpu = GetCpu()
	var disk, _ = GetDisk()
	var mem, _ = GetMemory()

	fmt.Println("[===============] HEARTBEAT [===============]")
	fmt.Println("[*] CPU Stats: ")
	fmt.Println("[+] User: ", cpu.User)
	fmt.Println("[+] Count: ", cpu.CpuCount)
	fmt.Println("[+] Total: ", cpu.Total)
	fmt.Println("[+] Idle: ", cpu.Idle)
	fmt.Println("[+] System", cpu.System)
	fmt.Println("")
	fmt.Println("[*] Memory stats: ")
	fmt.Println("[+] Cached: ", mem.Cached)
	fmt.Println("[+] Available: ", mem.Available)
	fmt.Println("[+] Active: ", mem.Active)
	fmt.Println("[+] Inactive: ", mem.Inactive)
	fmt.Println("[+] Free: ", mem.Free)
	fmt.Println("[+] Buffers: ", mem.Buffers)
}
