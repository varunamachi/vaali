package vapp

import (
	"errors"
	"net/http"

	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vsec"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vnet"
	"gopkg.in/mgo.v2/bson"
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
	}
}

func getEvents(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Fetch events")
	var events []*vlog.Event
	offset, limit, has := vnet.GetOffsetLimit(ctx)
	filter := make(bson.M)
	err = vnet.LoadJSONFromArgs(ctx, "filter", &filter)
	if err == nil && has {
		events, err = GetEvents(offset, limit, filter)
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
	err = ctx.JSON(status, &vnet.Result{
		Status: status,
		Op:     "fetch events",
		Msg:    msg,
		OK:     err == nil,
		Data:   events,
		Err:    err,
	})
	return vlog.LogError("App:Events", err)
}
