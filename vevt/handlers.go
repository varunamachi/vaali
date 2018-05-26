package vevt

import (
	"errors"
	"net/http"

	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vmgo"

	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vsec"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vnet"
)

func getEndpoints() []*vnet.Endpoint {
	return []*vnet.Endpoint{
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "admin/event",
			Access:   vsec.Admin,
			Category: "administration",
			Func:     getEvents,
			Comment:  "Fetch all the events",
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "admin/event/filterModel",
			Access:   vsec.Admin,
			Category: "administration",
			Func:     getEventFilterModel,
			Comment:  "Fetch event filter model",
		},
	}
}

func getEvents(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Fetch events")
	var events []*vlog.Event
	var total int
	offset, limit, has := vnet.GetOffsetLimit(ctx)
	var filter vmgo.Filter
	err = vnet.LoadJSONFromArgs(ctx, "filter", &filter)
	if err == nil && has {
		total, events, err = GetEvents(offset, limit, &filter)
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

func getEventFilterModel(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Get Event Filter Model")
	var efm EventFilterModel
	efm, err = GetEventFilterModel()
	if err != nil {
		msg = "Failed to retrieve event filter model"
		status = http.StatusInternalServerError
	}
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "event_filter_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   efm,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("Sec:Hdl", err)
}
