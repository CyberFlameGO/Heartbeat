package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {

}

//get the cpu usage and the idle time
func getCPUUsage() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")

	if err != nil {
		fmt.Println("Could not read /proc/stat found error: ", err)
		return
	}

	lines := strings.Split(string(contents), "\n")

	for line := range lines {
		fields := strings.Fields(string(line))
		if fields[0] == "cpu" {
			fieldCount := len(fields)

			for i := 1; i < fieldCount; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)

				if err != nil {
					total += val

					if i == 4 {
						idle = val
					}
				}

				return
			}
		}
	}

	return
}
