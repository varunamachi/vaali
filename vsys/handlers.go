package vsys

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vnet"
	"github.com/varunamachi/vaali/vsec"
)

//GetEndpoints - get REST endpoints related to operating and runtime systems
func GetEndpoints() []*vnet.Endpoint {
	return []*vnet.Endpoint{
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "system/stats",
			Access:   vsec.Admin,
			Category: "monitoring",
			Func:     getSysStats,
			Comment:  "Fetch system statistics",
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "system/info",
			Access:   vsec.Admin,
			Category: "monitoring",
			Func:     getSysInfo,
			Comment:  "Fetch disk usage info",
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "system/stats/cpu",
			Access:   vsec.Admin,
			Category: "monitoring",
			Func:     getCPUUsage,
			Comment:  "Fetch CPU usage info",
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "system/stats/memory",
			Access:   vsec.Admin,
			Category: "monitoring",
			Func:     getMemoryUsage,
			Comment:  "Fetch memory usage info",
		},
	}
}

func getSysStats(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Fetch System Stats")
	var stats *SysStat
	stats, err = GetSysStats()
	if err != nil {
		status = http.StatusInternalServerError
		msg = "Failed to fetch system stats"
	}
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "sys_stats_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   stats,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("App:Events", err)
}

func getSysInfo(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Fetch system info")
	var stats *SysInfo
	stats, err = GetSysInfo()
	if err != nil {
		status = http.StatusInternalServerError
		msg = "Failed to fetch system info"
	}
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "sys_info_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   stats,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("Sys:Stats", err)
}

func getCPUUsage(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Fetch cpu usage info")
	var stats *CPUStats
	stats, err = GetCPUStats()
	if err != nil {
		status = http.StatusInternalServerError
		msg = "Failed to fetch cpu usage info"
	}
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "stats_cpu_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   stats,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("App:Events", err)
}

func getMemoryUsage(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Fetch memory usage info")
	var stats *MemoryStats
	stats, err = GetMemStats()
	if err != nil {
		status = http.StatusInternalServerError
		msg = "Failed to fetch memory usage info"
	}
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "stats_cpu_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   stats,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("App:Events", err)
}
