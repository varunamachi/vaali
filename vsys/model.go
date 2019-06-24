package vsys

//DiskStats - disk statistics
type DiskStats struct {
	Path   string `json:"path"`
	Fstype string `json:"fstype"`
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
	Used   uint64 `json:"used"`
}

//MemoryStats - memory statistics
type MemoryStats struct {
	Total uint64 `json:"total"`
	Usage uint64 `json:"usage"`
	Free  uint64 `json:"free"`
}

//CPUStats - Statistics about the CPU
type CPUStats struct {
	NumPhysical int       `json:"numPhysical"`
	NumLogical  int       `json:"numLogical"`
	Usage       []float64 `json:"usage"`
}

//SysStat - system stats for the server
type SysStat struct {
	CPUStats    CPUStats    `json:"cpuUsage"`
	DiskStats   []DiskStats `json:"diskStats"`
	MemoryStats MemoryStats `json:"memoryStats"`
}
