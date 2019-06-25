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
			Func:     getSystemStats,
			Comment:  "Fetch system statistics",
		},
	}
}

func getSystemStats(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Fetch System Stats")
	stats, err := GetSystemStats()
	if err != nil {
		status = http.StatusInternalServerError
		msg = "Failed to fetch system stats"
	}
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "events_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   stats,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("App:Events", err)
}
