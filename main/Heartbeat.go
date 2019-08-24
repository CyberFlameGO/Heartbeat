package main

import (
	"fmt"
	"runtime"
)

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("[X] Could not start Heartbeat due to not being on linux, aborting!")
		runtime.StopTrace()
		return
	}
}
