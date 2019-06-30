package vsys

import (
	"time"

	"github.com/varunamachi/vaali/vlog"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/gopsutil/cpu"
)

//GetSysStats - get system statistics at the time of the call
func GetSysStats() (sysStats *SysStat, err error) {
	defer vlog.LogErrorX("Sys:Stat", "Failed to retrieve system stats", err)
	sysStats = &SysStat{}
	if sysStats.CPUStats, err = GetCPUStats(); err != nil {
		sysStats = nil
		return
	}
	if sysStats.MemoryStats, err = GetMemStats(); err != nil {
		sysStats = nil
		return
	}
	return
}

//GetSysInfo - gets system info
func GetSysInfo() (sysInfo *SysInfo, err error) {
	defer vlog.LogErrorX("Sys:Stat", "Failed to retrieve system stats", err)
	sysInfo = &SysInfo{}
	if sysInfo.DiskInfo, err = GetDiskInfo(); err != nil {
		sysInfo = nil
		return
	}

	if sysInfo.CPUInfo, err = GetCPUInfo(); err != nil {
		sysInfo = nil
		return
	}
	return
}

//GetCPUStats - get system CPU stats, number of cores, usage...
func GetCPUStats() (cpuStats *CPUStats, err error) {
	defer vlog.LogErrorX("Sys:Stats", "Failed to retrieve CPU stats", err)
	cpuStats = &CPUStats{}
	//Per core usage:
	// if cpuStats.Usage, err = cpu.Percent(1*time.Millisecond, true); err != nil {
	if cpuStats.Usage, err = cpu.Percent(0, true); err != nil {
		return
	}

	//All core usage, combined
	var combinedUsage []float64
	if combinedUsage, err = cpu.Percent(
		1*time.Millisecond, false); err != nil || len(combinedUsage) < 1 {
		return
	}
	cpuStats.CombinedUsage = combinedUsage[0]
	return
}

//GetMemStats - get memory stats
func GetMemStats() (memStats *MemoryStats, err error) {
	defer vlog.LogErrorX("Sys:Stats", "Failed to retrieve memory stats", err)
	mem, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	memStats = &MemoryStats{}
	memStats.Total = mem.Total
	memStats.Free = mem.Free
	memStats.Used = mem.Used
	return
}

//GetDiskInfo - get stats for all the disks in the system
func GetDiskInfo() (diskInfo []*DiskInfo, err error) {
	defer vlog.LogErrorX("Sys:Stats", "Failed to retrieve disk stats", err)
	partitions, err := disk.Partitions(false)
	if err != nil {
		return
	}
	diskInfo = make([]*DiskInfo, 0, len(partitions))
	for _, ptn := range partitions {
		usageStat, err := disk.Usage(ptn.Mountpoint)
		if err != nil {
			break
		}
		diskInfo = append(diskInfo, &DiskInfo{
			Path:   usageStat.Path,
			Fstype: usageStat.Fstype,
			Total:  usageStat.Total,
			Free:   usageStat.Free,
			Used:   usageStat.Used,
		})
	}
	return
}

//GetCPUInfo - gives information about the CPU/CPUs of ther server
func GetCPUInfo() (cpuInfo *CPUInfo, err error) {
	defer vlog.LogErrorX("Sys:Stats", "Failed to retrieve CPU info", err)
	cpuInfo = &CPUInfo{}

	//Num physical cores:
	if cpuInfo.NumPhysical, err = cpu.Counts(false); err != nil {
		return
	}
	//Num threads:
	if cpuInfo.NumLogical, err = cpu.Counts(true); err != nil {
		return
	}
	return
}
