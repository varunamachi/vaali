package vsec

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vnet"
	"github.com/varunamachi/vaali/vsec"
)

func createUser(ctx echo.Context) (err error) {
	op := "Create User"
	status, msg := vnet.DefMS(op)
	var user vsec.User
	err = ctx.Bind(&user)
	if err == nil {
		err = CreateUser(&user)
		if err != nil {
			msg = "Failed to create user in database"
			status = http.StatusInternalServerError
		}
	} else {
		status = http.StatusBadRequest
		msg = "User information given is malformed"
	}
	ctx.JSON(status, &vnet.Result{
		Op:   op,
		Msg:  msg,
		OK:   err == nil,
		Data: nil,
		Err:  err,
	})
	vlog.LogEvent(op, vnet.GetUserID(ctx), err == nil, err, nil)
	return vlog.LogError("Sec:Hdl", err)
}

func updateUser(ctx echo.Context) (err error) {
	return vlog.LogError("Sec:Hdl", err)
}

func deleteUser(ctx echo.Context) (err error) {
	return vlog.LogError("Sec:Hdl", err)
}

func getUser(ctx echo.Context) (err error) {
	return vlog.LogError("Sec:Hdl", err)
}

func getUsers(ctx echo.Context) (err error) {
	return vlog.LogError("Sec:Hdl", err)
}

func setPassword(ctx echo.Context) (err error) {
	return vlog.LogError("Sec:Hdl", err)
}

func resetPassword(ctx echo.Context) (err error) {
	return vlog.LogError("Sec:Hdl", err)
}
