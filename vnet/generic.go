package vnet

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vmgo"
)

func create(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Create", dtype)
	var data vmgo.StoredItem
	if len(dtype) != 0 {
		data, err = bind(ctx, dtype)
		if err == nil {
			data.SetCreationInfo(time.Now(), GetString(ctx, "userID"))
			err = vmgo.Create(dtype, data)
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
		Op:     dtype + "_create_",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("S:Entity", err)
}

func update(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Update", dtype)
	var data vmgo.StoredItem
	if len(dtype) != 0 {
		data, err = bind(ctx, dtype)
		if err == nil {
			data.SetModInfo(time.Now(), GetString(ctx, "userID"))
			err = vmgo.Update(dtype, bson.M{"_id": data.ID()}, data)
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
		Op:     dtype + "_update",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("S:Entity", err)
}

func delete(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Delete", dtype)
	id := ctx.Param("id")
	if len(dtype) != 0 {
		err = vmgo.Delete(dtype, bson.M{"_id": bson.ObjectId(id)})
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
		Op:     dtype + "_delete",
		Msg:    msg,
		OK:     err == nil,
		Data:   id,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("S:Entity", err)
}

func get(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Get", dtype)
	data := M{}
	id := ctx.Param("id")
	if len(dtype) != 0 {
		err = vmgo.Get(dtype, bson.M{"_id": bson.ObjectId(id)}, &data)
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
		Op:     dtype + "_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   data,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("S:Entity", err)
}

func getAll(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Get All", dtype)
	var data []*M
	if len(dtype) != 0 {
		offset, limit, has := GetOffsetLimit(ctx)
		var filter vmgo.Filter
		err = LoadJSONFromArgs(ctx, "filter", &filter)
		if has && err == nil {
			data = make([]*M, 0, limit)
			err = vmgo.GetAll(dtype, "-createdAt", offset, limit, &filter, data)
			if err != nil {
				msg = fmt.Sprintf("Failed to retrieve %s from database", dtype)
				status = http.StatusInternalServerError
			}
		} else {
			msg = "Invalid offset, limit or filter given"
			status = http.StatusBadRequest
		}
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
	}
	err = SendAndAuditOnErr(ctx, &Result{
		Status: status,
		Op:     dtype + "_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   data,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("S:Entity", err)
}

func count(ctx echo.Context) (err error) {
	//@TODO - handle filters
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Get All", dtype)
	count := 0
	if len(dtype) != 0 {
		if err == nil {
			var filter vmgo.Filter
			err = LoadJSONFromArgs(ctx, "filter", &filter)
			count, err = vmgo.Count(dtype, &filter)
			if err != nil {
				msg = fmt.Sprintf("Failed to retrieve %s from database", dtype)
				status = http.StatusInternalServerError
			}
		} else {
			msg = fmt.Sprintf("Failed to decode filter for '%s'", dtype)
			status = http.StatusInternalServerError
		}
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
	}
	err = SendAndAuditOnErr(ctx, &Result{
		Status: status,
		Op:     dtype + "_count",
		Msg:    msg,
		OK:     err == nil,
		Data:   count,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("S:Entity", err)
}

func getFilterValues(ctx echo.Context) (err error) {
	dtype := ctx.Param("dataType")
	status, msg := DefaultSM("Filter Values of", dtype)
	var fdesc []*vmgo.FilterSpec
	var values bson.M
	if len(dtype) != 0 {
		err = LoadJSONFromArgs(ctx, "fdesc", &fdesc)
		if err == nil {
			values, err = vmgo.GetFilterValues(dtype, fdesc)
		} else {
			msg = "Failed to load filter description from URL"
			status = http.StatusBadRequest
		}
	} else {
		msg = "Invalid empty data type given"
		status = http.StatusBadRequest
		err = errors.New(msg)
	}
	err = SendAndAuditOnErr(ctx, &Result{
		Status: status,
		Op:     dtype + "_filter_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data:   values,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("S:Entity", err)
}

func bind(ctx echo.Context, dataType string) (
	data vmgo.StoredItem, err error) {
	data = vmgo.Instance(dataType)
	if data == nil {
		err = ctx.Bind(data)
	} else {
		err = fmt.Errorf("Could not find factory function for data type %s",
			dataType)
	}
	return data, err
}
