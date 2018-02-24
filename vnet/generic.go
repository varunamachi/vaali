package vnet

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vdb"
	"github.com/varunamachi/vaali/vlog"
)

func create(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Create", dtype)
	data := bson.M{}
	if len(dtype) != 0 {
		err = ctx.Bind(&data)
		if err == nil {
			data["createdAt"] = time.Now()
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
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
	}
	AuditedSendX(ctx, &data, &Result{
		Status: status,
		Op:     "create_" + dtype,
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    err,
	})
	return vlog.LogError("S:Entity", err)
}

func update(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Update", dtype)
	data := M{}
	if len(dtype) != 0 {
		err = ctx.Bind(&data)
		if err == nil {
			id, _ := data["_id"].(string)
			data["modifiedAt"] = time.Now()
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
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
	}
	AuditedSendX(ctx, &data, &Result{
		Status: status,
		Op:     "update_" + dtype,
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    err,
	})
	return vlog.LogError("S:Entity", err)
}

func delete(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Delete", dtype)
	id := ctx.Param("id")
	if len(dtype) != 0 {
		err = vdb.Delete(dtype, bson.M{"_id": bson.ObjectId(id)})
		if err != nil {
			msg = fmt.Sprintf("Failed to delete %s from database", dtype)
			status = http.StatusInternalServerError
		}
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
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

func get(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Get", dtype)
	data := M{}
	id := ctx.Param("id")
	if len(dtype) != 0 {
		err = vdb.Get(dtype, bson.M{"_id": bson.ObjectId(id)}, &data)
		if err != nil {
			msg = fmt.Sprintf(
				"Failed to retrieve %s from database, entity with ID %s",
				dtype,
				id)
			status = http.StatusInternalServerError
		}
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
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

func getAll(ctx echo.Context) (err error) {
	//@TODO - handle filters
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Get All", dtype)
	var data []*M
	if len(dtype) != 0 {
		offset, limit, has := GetOffsetLimit(ctx)
		if has {
			data = make([]*M, 0, limit)
			err = vdb.GetAll(dtype, "-createdAt", offset, limit, data)
			if err != nil {
				msg = fmt.Sprintf("Failed to retrieve %s from database", dtype)
				status = http.StatusInternalServerError
			}
		} else {
			msg = "Invalid offset and limit given"
			status = http.StatusBadRequest
		}
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
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

func count(ctx echo.Context) (err error) {
	//@TODO - handle filters
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Get All", dtype)
	count := 0
	if len(dtype) != 0 {
		count, err = vdb.Count(dtype)
		if err != nil {
			msg = fmt.Sprintf("Failed to retrieve %s from database", dtype)
			status = http.StatusInternalServerError
		}
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
	}
	err = SendAndAuditOnErr(ctx, &Result{
		Status: status,
		Op:     "get_" + dtype,
		Msg:    msg,
		OK:     err == nil,
		Data:   count,
		Err:    err,
	})
	return vlog.LogError("S:Entity", err)
}
