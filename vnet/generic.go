package vnet

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vdb"
	"github.com/varunamachi/vaali/vlog"
)

//CreateHandler - creates a generic POST handler for a data type
func CreateHandler(dtype, colnName string) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		status, msg := DefaultSM("Create", dtype)
		data := M{}
		err = ctx.Bind(&data)
		if err == nil {
			err = vdb.Create(colnName, data)
			if err != nil {
				msg = fmt.Sprintf("Failed to create %s in database", dtype)
				status = http.StatusInternalServerError
			}
		} else {
			msg = fmt.Sprintf(
				"Failed to retrieve %s information from the request", dtype)
			status = http.StatusBadRequest
		}
		err = AuditedSendX(ctx, &data, &Result{
			Status: status,
			Op:     "create_" + dtype,
			Msg:    msg,
			OK:     err == nil,
			Data:   nil,
			Err:    err,
		})
		return vlog.LogError("S:Entity", err)
	}
}

//CreateHandlerX - creates a generic POST handler for a data type, expects a
//custom logic for binding request body to a data struct
func CreateHandlerX(dtype, colnName string, bind BinderFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		status, msg := DefaultSM("Create", dtype)
		var data interface{}
		data, err = bind(ctx)
		if err == nil {
			err = vdb.Create(colnName, data)
			if err != nil {
				msg = fmt.Sprintf("Failed to create %s in database", dtype)
				status = http.StatusInternalServerError
			}
		} else {
			msg = fmt.Sprintf(
				"Failed to retrieve %s information from the request", dtype)
			status = http.StatusBadRequest
		}
		err = AuditedSendX(ctx, &data, &Result{
			Status: status,
			Op:     "create_" + dtype,
			Msg:    msg,
			OK:     err == nil,
			Data:   nil,
			Err:    err,
		})
		return vlog.LogError("S:Entity", err)
	}
}
