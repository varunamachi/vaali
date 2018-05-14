package vuman

import (
	"errors"
	"net/http"
	"time"

	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vmgo"

	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vnet"
	"github.com/varunamachi/vaali/vsec"
)

func updateUserInfo(user *vsec.User) (err error) {
	if len(user.ID) == 0 {
		// @TODO - store hash of user ID
		user.ID = vcmn.Hash(user.Email)
	} else {
		user.ID = vcmn.Hash(user.ID)
	}
	user.VerID = uuid.NewV4().String()
	user.Created = time.Now()
	user.State = vsec.Disabled
	// @TODO create a key retrieving strategy -- local | remote etc
	var emailKey string
	err = vcmn.GetConfig("emailKey", &emailKey)
	if err == nil {
		user.Email, err = vcmn.EncryptStr(emailKey, user.Email)
	}
	return err
}

func createUser(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Create User")
	var user vsec.User
	err = ctx.Bind(&user)
	if err == nil {
		user.Props = bson.M{
			"admin-created": true,
		}
		updateUserInfo(&user)
		err = CreateUser(&user)
		if err != nil {
			msg = "Failed to create user in database"
			status = http.StatusInternalServerError
		} else {
			err = sendVerificationMail(&user)
			// fmt.Println(getVerificationLink(&user))
			if err != nil {
				msg = "Failed to send verification email"
				status = http.StatusInternalServerError
			}
		}
	} else {
		status = http.StatusBadRequest
		msg = "User information given is malformed"
	}
	err = vnet.AuditedSendX(ctx, user, &vnet.Result{
		Status: status,
		Op:     "user_create",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("Sec:Hdl", err)
}

func registerUser(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Register User")
	// var user vsec.User
	upw := struct {
		User     vsec.User `json:"user"`
		Password string    `json:"password"`
	}{}
	err = ctx.Bind(&upw)
	if err == nil {
		upw.User.Auth = vsec.Normal
		updateUserInfo(&upw.User)
		err = CreateUser(&upw.User)
		if err != nil {
			msg = "Failed to register user in database"
			status = http.StatusInternalServerError
		} else {

			err = SetPassword(upw.User.ID, upw.Password)
			if err != nil {
				msg = "Failed to set password"
				status = http.StatusInternalServerError
			} else {
				err = sendVerificationMail(&upw.User)
				if err != nil {
					msg = "Failed to send verification email"
					status = http.StatusInternalServerError
				}
			}
		}
	} else {
		status = http.StatusBadRequest
		msg = "User information given is malformed"
	}
	err = vnet.AuditedSendX(ctx, upw.User, &vnet.Result{
		Status: status,
		Op:     "user_register",
		Msg:    msg,
		OK:     err == nil,
		Data: vlog.M{
			"user": upw.User,
		},
		Err: vcmn.ErrString(err),
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
	err = vnet.AuditedSendX(ctx, user, &vnet.Result{
		Status: status,
		Op:     "user_update",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("Sec:Hdl", err)
}

func deleteUser(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Delete User")
	userID := ctx.Param("userID")
	var user *vsec.User
	user, err = GetUser(userID)
	if err == nil {
		curID := vnet.GetString(ctx, "userID")
		if userID == curID {
			msg = "Can not delete own user account"
			status = http.StatusBadRequest
		} else if user.Auth == vsec.Super {
			msg = "Super account can not be deleted from web interface"
			status = http.StatusBadRequest
			err = errors.New(msg)
		} else {
			err = DeleteUser(userID)
			if err != nil {
				msg = "Failed to delete user from database"
				status = http.StatusInternalServerError
			}
		}
	} else {
		msg = "Invalid user ID is given for deletion"
		status = http.StatusBadRequest
		err = errors.New(msg)
	}
	err = vnet.AuditedSend(ctx, &vnet.Result{
		Status: status,
		Op:     "user_remove",
		Msg:    msg,
		OK:     err == nil,
		Data: vlog.M{
			"id":   userID,
			"user": user,
		},
		Err: vcmn.ErrString(err),
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
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "user_get",
		Msg:    msg,
		OK:     err == nil,
		Data:   user,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("Sec:Hdl", err)
}

func getUsers(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Get Users")
	offset, limit, has := vnet.GetOffsetLimit(ctx)
	// us := vcmn.GetFirstValidStr(ctx.Param("status"), string(vsec.Active))
	var users []*vsec.User
	var total int
	var filter vmgo.Filter
	err = vnet.LoadJSONFromArgs(ctx, "filter", &filter)
	if has && err == nil {
		total, users, err = GetUsers(offset, limit, &filter)
		if err != nil {
			msg = "Failed to retrieve user info from database"
			status = http.StatusInternalServerError
		}
	} else if err != nil {
		msg = "Failed to decode filter"
		status = http.StatusBadRequest
	} else {
		msg = "Could not retrieve user list, offset/limit not found"
		status = http.StatusBadRequest
		err = errors.New(msg)
	}
	err = vnet.SendAndAuditOnErr(ctx, &vnet.Result{
		Status: status,
		Op:     "user_multi_fetch",
		Msg:    msg,
		OK:     err == nil,
		Data: vmgo.CountList{
			TotalCount: total,
			Data:       users,
		},
		Err: vcmn.ErrString(err),
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
	err = vnet.AuditedSendX(ctx, vcmn.Hash(userID), &vnet.Result{
		Status: status,
		Op:     "user_password_set",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("Sec:Hdl", err)
}

func resetPassword(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Set Password")
	pinfo := make(map[string]string)
	err = ctx.Bind(&pinfo)
	userID := vnet.GetString(ctx, "userID")
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
	err = vnet.AuditedSendX(ctx, vcmn.Hash(userID), &vnet.Result{
		Status: status,
		Op:     "user_password_reset",
		Msg:    msg,
		OK:     err == nil,
		Data:   nil,
		Err:    vcmn.ErrString(err),
	})
	return vlog.LogError("Sec:Hdl", err)
}

// func verifyUser(ctx echo.Context) (err error) {
// 	status, msg := vnet.DefMS("User verification")
// 	userID := ctx.Param("userID")
// 	verID := ctx.Param("verID")
// 	if len(userID) != 0 && len(verID) != 0 {
// 		err = VerifyUser(userID, verID)
// 		if err != nil {
// 			msg = "Failed to verify user"
// 			status = http.StatusInternalServerError
// 		}
// 	} else {
// 		status = http.StatusBadRequest
// 		msg = "User information given is malformed"
// 	}
// 	vnet.AuditedSendX(ctx, userID, &vnet.Result{
// 		Status: status,
// 		Op:     "user_verify",
// 		Msg:    msg,
// 		OK:     err == nil,
// 		Data: vlog.M{
// 			"userID": userID,
// 			"verID":  verID,
// 		},
// 		Err: err,
// 	})
// 	return vlog.LogError("Sec:Hdl", err)
// }

func verify(ctx echo.Context) (err error) {
	status, msg := vnet.DefMS("Create Password")
	params := make(map[string]string)
	userID := ctx.Param("userID")
	verID := ctx.Param("verID")
	err = ctx.Bind(&params)
	if len(userID) > 0 && len(verID) > 0 && err == nil {
		err = VerifyUser(userID, verID)
		if err == nil {
			err = SetPassword(userID, params["password"])
			if err != nil {
				msg = "Failed to set password"
				status = http.StatusInternalServerError
			}
		} else {
			msg = "Failed to verify user"
			status = http.StatusInternalServerError
		}
	} else {
		status = http.StatusBadRequest
		msg = "Invalid information provided for creating password"
	}
	ctx.Set("userName", "N/A")
	hash := vcmn.Hash(userID)
	err = vnet.AuditedSendX(ctx, hash, &vnet.Result{
		Status: status,
		Op:     "user_account_verify",
		Msg:    msg,
		OK:     err == nil,
		Data: vlog.M{
			"userID":         hash,
			"verificationID": verID,
		},
		Err: vcmn.ErrString(err),
	})
	return vlog.LogError("Sec:Hdl", err)
}
