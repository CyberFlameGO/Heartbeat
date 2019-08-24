package main

import (
	"fmt"
	"os"
)

type Stats struct {
	User, Nice, System, Idle, Iowait, Irq, Softirq, Steal, Guest, GuestNice, Total uint64
	CpuCount, StatCount                                                            uint64
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

	//todo return getting the stats
}

func main() {

}
