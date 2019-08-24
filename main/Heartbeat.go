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
	} else {

		var cpu, _ = GetCpuStats()
		var memory, _ = GetMemory()
		//FIXME: var disk, _ = GetDiskStats()

		fmt.Println("{==== HEARTBEAT ====}")
		fmt.Println("[*] CPU Stats: ")
		fmt.Println("[+] User: ", cpu.User)
		fmt.Println("[+] Count: ", cpu.CPUCount)
		fmt.Println("[+] Total: ", cpu.Total)
		fmt.Println("[+] Idle: ", cpu.Idle)
		fmt.Println("[+] System", cpu.System)
		fmt.Println("")
		fmt.Println("[*] Memory stats: ")
		fmt.Println("[+] Cached: ", memory.Cached)
		fmt.Println("[+] Available: ", memory.Available)
		fmt.Println("[+] Active: ", memory.Active)
		fmt.Println("[+] Inactive: ", memory.Inactive)
		fmt.Println("[+] Free: ", memory.Free)
		fmt.Println("[+] Buffers: ", memory.Buffers)
		fmt.Println("[+] Total: ", memory.Total)
		fmt.Println("[+] Used: ", memory.Used)
		fmt.Println("")
		fmt.Println("[*] Disk stats: ")
		fmt.Println("[X] I broke this...ops")
		fmt.Println("{==== HEARTBEAT ====}")
	}
}
