package vapp

import (
	"errors"
	"net/http"

	"github.com/varunamachi/vaali/vevt"

	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vmgo"
	"github.com/varunamachi/vaali/vnet"

	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vsec"

	"github.com/labstack/echo"
)

func getEndpoints() []*vnet.Endpoint {
	return []*vnet.Endpoint{
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "event",
			Access:   vsec.Admin,
			Category: "administration",
			Func:     getEvents,
			Comment:  "Fetch all the events",
		},
	}
}

func getEvents(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Fetch events")
	var events []*vevt.Event
	var total int
	offset, limit, has := vnet.GetOffsetLimit(ctx)
	var filter vmgo.Filter
	err = vnet.LoadJSONFromArgs(ctx, "filter", &filter)
	if err == nil && has {
		total, events, err = vevt.GetAuditor().GetEvents(
			offset, limit, &filter)
		if err != nil {
			msg = "Could not retrieve event info from database"
			status = http.StatusInternalServerError
		}
	} else {
		if err == nil {
			err = errors.New("Could not get Offset and Limit arguments")
		}
		msg = "Could not find required parameter"
		status = http.StatusBadRequest
	}
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "events_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data: vmgo.CountList{
			TotalCount: total,
			Data:       events,
		},
		Err: vcmn.ErrString(err),
	})
	return vlog.LogError("App:Events", err)
}
