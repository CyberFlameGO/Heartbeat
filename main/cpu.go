package main

type CPUStats struct {
	User, System, Idle, Total uint64
	CPUCount, StatCount       uint64
}

type CPUInfo struct {
	name string
	ptr  *uint64
}
