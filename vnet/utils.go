package vnet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/varunamachi/vaali/vcmn"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vlog"
)

//GetOffsetLimit - retrieves offset and limit as integers provided in URL as
//query parameters. These parameters should have name - offset and limit
//respectively
func GetOffsetLimit(ctx echo.Context) (offset, limit int, has bool) {
	has = false
	offset = 0
	limit = 0
	strOffset := ctx.QueryParam("offset")
	strLimit := ctx.QueryParam("limit")
	if len(strOffset) == 0 || len(strLimit) == 0 {

		has = false
		return
	}
	var err error
	offset, err = strconv.Atoi(strOffset)
	if err != nil {
		offset = 0
		return
	}
	limit, err = strconv.Atoi(strLimit)
	if err != nil {
		limit = 0
		return
	}
	has = true
	return
}

//DefMS - gives default message and status, used for convenience
func DefMS(oprn string) (int, string) {
	return http.StatusOK, oprn + " - successful"
}

//DefaultSM - default status and message
func DefaultSM(opern, name string) (int, string) {
	return http.StatusOK, fmt.Sprintf("%s %s - successful", opern, name)
}

// //GetUserID - retrieves user ID from context
// func GetUserID_(ctx echo.Context) string {
// 	ui := ctx.Get("userID")
// 	userID, ok := ui.(string)
// 	if ok {
// 		return userID
// 	}
// 	return ""
// }

//GetString - retrieves property with key from context
func GetString(ctx echo.Context, key string) (value string) {
	ui := ctx.Get(key)
	userID, ok := ui.(string)
	if ok {
		return userID
	}
	return ""
}

//AuditedSend - sends result as JSON while logging it as event. The event data
//is same as the data present in the result
func AuditedSend(ctx echo.Context, res *Result) (err error) {
	err = ctx.JSON(res.Status, res)
	vlog.LogEvent(
		res.Op,
		GetString(ctx, "userID"),
		GetString(ctx, "userName"),
		res.OK,
		res.Err,
		res.Data)
	return err
}

//AuditedSendX - sends result as JSON while logging it as event. This method
//logs event data which is seperate from result data
func AuditedSendX(ctx echo.Context, data interface{}, res *Result) (err error) {
	err = ctx.JSON(res.Status, res)
	vlog.LogEvent(
		res.Op,
		GetString(ctx, "userID"),
		GetString(ctx, "userName"),
		res.OK,
		res.Err,
		data)
	return err
}

//SendAndAuditOnErr - sends the result to client and puts an audit record in
//audit log if the result is error OR sending failed
func SendAndAuditOnErr(ctx echo.Context, res *Result) (err error) {
	fmt.Println("one")
	err = ctx.JSON(res.Status, res)
	if res.Err != nil || err != nil {
		fmt.Println("two")
		vlog.LogEvent(
			res.Op,
			GetString(ctx, "userID"),
			GetString(ctx, "userName"),
			false,
			vcmn.FirstValid(res.Err, err),
			res.Data)
	}
	return err
}

//LoadJSONFromArgs - decodes argument identified by 'param' to JSON and
//unmarshals it into the given 'out' structure
func LoadJSONFromArgs(ctx echo.Context, param string, out interface{}) (
	err error) {
	val := ctx.QueryParam(param)
	if len(val) != 0 {
		var decoded string
		decoded, err = url.PathUnescape(val)
		if err == nil {
			err = json.Unmarshal([]byte(decoded), out)
		}
	}
	return vlog.LogError("Net:Utils", err)
}

//Merge - merges multple endpoint arrays
func Merge(epss ...[]*Endpoint) (eps []*Endpoint) {
	eps = make([]*Endpoint, 0, 100)
	for _, es := range epss {
		eps = append(eps, es...)
	}
	return eps
}
