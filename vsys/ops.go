package vsys

import (
	"time"

	"github.com/varunamachi/vaali/vlog"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/gopsutil/cpu"
)

//GetSystemStats - get system stats for the server
func GetSystemStats() (sysStats *SysStat, err error) {
	defer func() {
		vlog.LogErrorX("Sys:Stat", "Failed to retrieve system stats", err)
	}()
	sysStats = &SysStat{
		DiskStats: make([]DiskStats, 0, 20),
	}

	if sysStats.CPUStats.NumPhysical, err = cpu.Counts(false); err != nil {
		return
	}
	if sysStats.CPUStats.NumLogical, err = cpu.Counts(true); err != nil {
		return
	}
	if sysStats.CPUStats.Usage, err = cpu.Percent(
		1*time.Millisecond, true); err != nil {
		return
	}
	mem, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	sysStats.MemoryStats.Total = mem.Total
	sysStats.MemoryStats.Free = mem.Free
	sysStats.MemoryStats.Usage = mem.Used

	partitions, err := disk.Partitions(false)
	if err != nil {
		return
	}
	for _, ptn := range partitions {
		usageStat, err := disk.Usage(ptn.Mountpoint)
		if err != nil {
			return sysStats, err
		}
		sysStats.DiskStats = append(sysStats.DiskStats, DiskStats{
			Path:   usageStat.Path,
			Fstype: usageStat.Fstype,
			Total:  usageStat.Total,
			Free:   usageStat.Free,
			Used:   usageStat.Used,
		})
	}
	return
}
