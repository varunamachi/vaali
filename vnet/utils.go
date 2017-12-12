package vnet

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
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

//GetUserID - retrieves user ID from context
func GetUserID(ctx echo.Context) string {
	ui := ctx.Get("userID")
	userID, ok := ui.(string)
	if ok {
		return userID
	}
	return ""
}
