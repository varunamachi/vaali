package vsys

import (
	_ "github.com/shirou/gopsutil" //For now
)

//SysStat - system stats for the server
type SysStat struct {
	CPUUsage []float64 `json:"cpuUsage"`
}
