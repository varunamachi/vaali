package vsys

//MemoryStats - memory statistics
type MemoryStats struct {
	Total uint64 `json:"total"`
	Usage uint64 `json:"usage"`
	Free  uint64 `json:"free"`
}

//CPUStats - Statistics about the CPU
type CPUStats struct {
	Usage         []float64 `json:"usage"`
	CombinedUsage float64   `json:"combinedUsage"`
}

//SysStat - system stats for the server
type SysStat struct {
	CPUStats    *CPUStats    `json:"cpuUsage"`
	MemoryStats *MemoryStats `json:"memoryStats"`
}

//DiskInfo - disk info and usage
type DiskInfo struct {
	Path   string `json:"path"`
	Fstype string `json:"fstype"`
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
	Used   uint64 `json:"used"`
}

//CPUInfo - cpu information
type CPUInfo struct {
	NumPhysical int `json:"numPhysical"`
	NumLogical  int `json:"numLogical"`
}

//SysInfo - system info, less frequently changing
type SysInfo struct {
	CPUInfo  *CPUInfo    `json:"cpuInfo"`
	DiskInfo []*DiskInfo `json:"diskInfo"`
}
