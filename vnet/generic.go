package vnet

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vdb"
	"github.com/varunamachi/vaali/vlog"
)

//MakeCreateHandler - creates a generic POST handler for a data type
func MakeCreateHandler(dtype string) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		status, msg := DefaultSM("Create", dtype)
		data := M{}
		err = ctx.Bind(&data)
		if err == nil {
			err = vdb.Create(dtype, data)
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

//MakeCreateHandlerX - creates a generic POST handler for a data type, expects a
//custom logic for binding request body to a data struct
func MakeCreateHandlerX(dtype string, bind BinderFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		status, msg := DefaultSM("Create", dtype)
		var data interface{}
		data, err = bind(ctx)
		if err == nil {
			err = vdb.Create(dtype, data)
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

//MakeUpdateHandler - creates a generic PUT handler for a updating data.
//Object ID is used for matching
func MakeUpdateHandler(dtype string) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		status, msg := DefaultSM("Update", dtype)
		data := M{}
		err = ctx.Bind(&data)
		if err == nil {
			id, _ := data["_id"].(string)
			err = vdb.Update(dtype, bson.M{"_id": bson.ObjectIdHex(id)}, data)
			if err != nil {
				msg = fmt.Sprintf("Failed to update %s in database", dtype)
				status = http.StatusInternalServerError
			}
		} else {
			msg = fmt.Sprintf(
				"Failed to retrieve %s information from the request", dtype)
			status = http.StatusBadRequest
		}
		err = AuditedSendX(ctx, &data, &Result{
			Status: status,
			Op:     "update_" + dtype,
			Msg:    msg,
			OK:     err == nil,
			Data:   nil,
			Err:    err,
		})
		return vlog.LogError("S:Entity", err)
	}
}

//MakeUpdateHandlerX - creates a generic PUT handler for a data, expects a
//custom logic for binding request body to a data struct. Object ID is used for
//matching
// func MakeUpdateHandlerX(dtype string,
// 	bind BinderFunc) echo.HandlerFunc {
// 	return func(ctx echo.Context) (err error) {
// 		status, msg := DefaultSM("Update", dtype)
// 		var data interface{}
// 		data, err = bind(ctx)
// 		if err == nil {
// 			id, _ := data["_id"].(string)
// 			err = vdb.Update(dtype, bson.M{"_id": bson.ObjectIdHex(id)}, data)
// 			if err != nil {
// 				msg = fmt.Sprintf("Failed to update %s in database", dtype)
// 				status = http.StatusInternalServerError
// 			}
// 		} else {
// 			msg = fmt.Sprintf(
// 				"Failed to retrieve %s information from the request", dtype)
// 			status = http.StatusBadRequest
// 		}
// 		err = AuditedSendX(ctx, &data, &Result{
// 			Status: status,
// 			Op:     "update_" + dtype,
// 			Msg:    msg,
// 			OK:     err == nil,
// 			Data:   nil,
// 			Err:    err,
// 		})
// 		return vlog.LogError("S:Entity", err)
// 	}
// }

//MakeDeleteHandler - creates a delete handler function, uses a REST parameter
//named 'id' for getting the Object ID of the item to be deleted
func MakeDeleteHandler(dtype string) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		status, msg := DefaultSM("Delete", dtype)
		id := ctx.Param("id")
		err = vdb.Delete(dtype, bson.M{"_id": bson.ObjectId(id)})
		if err != nil {
			msg = fmt.Sprintf("Failed to delete %s from database", dtype)
			status = http.StatusInternalServerError
		}
		err = AuditedSend(ctx, &Result{
			Status: status,
			Op:     "delete_" + dtype,
			Msg:    msg,
			OK:     err == nil,
			Data:   id,
			Err:    err,
		})
		return vlog.LogError("S:Entity", err)
	}
}

//MakeGetHandler - creates a GET handler function, uses a REST parameter
//named 'id' for getting the Object ID of the item to be retrieved
func MakeGetHandler(dtype string) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		status, msg := DefaultSM("Get", dtype)
		data := M{}
		id := ctx.Param("id")
		err = vdb.Get(dtype, bson.M{"_id": bson.ObjectId(id)}, &data)
		if err != nil {
			msg = fmt.Sprintf(
				"Failed to retrieve %s from database, entity with ID %s",
				dtype,
				id)
			status = http.StatusInternalServerError
		}
		err = SendAndAuditOnErr(ctx, &Result{
			Status: status,
			Op:     "get_" + dtype,
			Msg:    msg,
			OK:     err == nil,
			Data:   data,
			Err:    err,
		})
		return vlog.LogError("S:Entity", err)
	}
}

//MakeGetAllHandler - makes get all handler with offset and limit,
//Uses sort field for sorting records on descending order
func MakeGetAllHandler(dtype, sortField string) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		status, msg := DefaultSM("Get All", dtype)
		offset, limit, has := GetOffsetLimit(ctx)
		var data []*M
		if has {
			data = make([]*M, 0, limit)
			err = vdb.GetAll(dtype, "-" + sortField, offset, limit, data)
			if err != nil {
				msg = fmt.Sprintf("Failed to retrieve %s from database", dtype)
				status = http.StatusInternalServerError
			}
		}
		err = SendAndAuditOnErr(ctx, &Result{
			Status: status,
			Op:     "get_" + dtype,
			Msg:    msg,
			OK:     err == nil,
			Data:   data,
			Err:    err,
		})
		return vlog.LogError("S:Entity", err)
	}
}
