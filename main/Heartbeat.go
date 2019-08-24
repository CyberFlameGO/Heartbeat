package main

import (
	"errors"
	"io/ioutil"
)

func main() {

}

//cpu usage

func getCPUUsage() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")

	if err != nil {
		errors.New("Could not read /proc/stat")
		err.Error()
		return
	}
}
