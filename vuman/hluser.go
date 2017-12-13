package vsec

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vnet"
	"github.com/varunamachi/vaali/vsec"
)

//GetUserEndpoints - gives REST endpoints for user manaagement
func GetUserEndpoints() (endpoints []*vnet.Endpoint) {
	endpoints = []*vnet.Endpoint{
		&vnet.Endpoint{
			Method:   echo.POST,
			URL:      "/uman/user",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     createUser,
		},
		&vnet.Endpoint{
			Method:   echo.PUT,
			URL:      "/uman/user",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     updateUser,
		},
		&vnet.Endpoint{
			Method:   echo.DELETE,
			URL:      "/uman/user/:userID",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     deleteUser,
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "/uman/user/:userID",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     getUser,
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "/uman/user",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     getUsers,
		},
		&vnet.Endpoint{
			Method:   echo.POST,
			URL:      "/uman/user/password",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     setPassword,
		},
		&vnet.Endpoint{
			Method:   echo.PUT,
			URL:      "/uman/user/password",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     resetPassword,
		},
	}
	return endpoints

}

func createUser(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Create User")
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
	vnet.AuditedSendX(ctx, user, &vnet.Result{
		Status: status,
		Op:     "user_create",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    err,
	})
	return vlog.LogError("Sec:Hdl", err)
}

func updateUser(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Update User")
	var user vsec.User
	err = ctx.Bind(&user)
	if err == nil {
		err = UpdateUser(&user)
		if err != nil {
			msg = "Failed to update user in database"
			status = http.StatusInternalServerError
		}
	} else {
		status = http.StatusBadRequest
		msg = "User information given is malformed"
	}
	vnet.AuditedSendX(ctx, user, &vnet.Result{
		Status: status,
		Op:     "user_update",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    err,
	})
	return vlog.LogError("Sec:Hdl", err)
}

func deleteUser(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Delete User")
	userID := ctx.Param("userID")
	if len(userID) == 0 {
		err = DeleteUser(userID)
		if err != nil {
			msg = "Failed to delete user from database"
			status = http.StatusInternalServerError
		}
	} else {
		msg = "Invalid user ID is given for deletion"
		status = http.StatusBadRequest
	}
	vnet.AuditedSend(ctx, &vnet.Result{
		Status: status,
		Op:     "user_remove",
		Msg:    msg,
		OK:     err == nil,
		Data:   userID,
		Err:    err,
	})
	return vlog.LogError("Sec:Hdl", err)
}

func getUser(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Get User")
	userID := ctx.Param("userID")
	var user *vsec.User
	if len(userID) == 0 {
		user, err = GetUser(userID)
		if err != nil {
			msg = "Failed to retrieve user info from database"
			status = http.StatusInternalServerError
		}
	} else {
		msg = "Invalid user ID is given for retrieval"
		status = http.StatusBadRequest
	}
	vnet.AuditedSend(ctx, &vnet.Result{
		Status: status,
		Op:     "user_get",
		Msg:    msg,
		OK:     err == nil,
		Data:   user,
		Err:    err,
	})
	return vlog.LogError("Sec:Hdl", err)
}

func getUsers(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Get Users")
	offset, limit, has := vnet.GetOffsetLimit(ctx)
	var users []*vsec.User
	if has {
		users, err = GetAllUsers(offset, limit)
		if err != nil {
			msg = "Failed to retrieve user info from database"
			status = http.StatusInternalServerError
		}
	} else {
		msg = "Could not retrieve user list, offset/limit not found"
		status = http.StatusBadRequest
	}
	vnet.AuditedSendX(ctx, nil, &vnet.Result{
		Status: status,
		Op:     "multi_user_get",
		Msg:    msg,
		OK:     err == nil,
		Data:   users,
		Err:    err,
	})
	return vlog.LogError("Sec:Hdl", err)
}

func setPassword(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Set Password")
	pinfo := make(map[string]string)
	err = ctx.Bind(&pinfo)
	userID, ok1 := pinfo["userID"]
	password, ok2 := pinfo["password"]
	if err == nil && ok1 && ok2 {
		err = SetPassword(userID, password)
		if err != nil {
			msg = "Failed to set password in database"
			status = http.StatusInternalServerError
		}
	} else {
		status = http.StatusBadRequest
		msg = "Password information given is invalid, cannot set"
	}
	vnet.AuditedSendX(ctx, userID, &vnet.Result{
		Status: status,
		Op:     "password_set",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    err,
	})
	return vlog.LogError("Sec:Hdl", err)
}

func resetPassword(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Set Password")
	pinfo := make(map[string]string)
	err = ctx.Bind(&pinfo)
	userID := vnet.GetUserID(ctx)
	oldPassword, ok2 := pinfo["oldPassword"]
	newPassword, ok3 := pinfo["newPassword"]
	if err == nil && ok2 && ok3 && len(userID) != 0 {
		err = ResetPassword(userID, oldPassword, newPassword)
		if err != nil {
			msg = "Failed to reset password in database"
			status = http.StatusInternalServerError
		}
	} else {
		status = http.StatusBadRequest
		msg = "Password information given is invalid, cannot reset"
	}
	vnet.AuditedSendX(ctx, userID, &vnet.Result{
		Status: status,
		Op:     "password_reset",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    err,
	})
	return vlog.LogError("Sec:Hdl", err)
}
