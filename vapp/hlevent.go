package vapp

import (
	"errors"
	"net/http"

	"github.com/varunamachi/vaali/vlog"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vnet"
	"gopkg.in/mgo.v2/bson"
)

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
	err = vnet.AuditedSend(ctx, &vnet.Result{
		Status: status,
		Op:     "fetch events",
		Msg:    msg,
		OK:     err == nil,
		Data:   events,
		Err:    err,
	})
	return vlog.LogError("App:Events", err)
}
