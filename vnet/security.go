package vnet

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/varunamachi/vaali/vcmn"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vsec"
)

//JWTUserInfo - container for retrieving user information from JWT
type JWTUserInfo struct {
	UserID   string
	UserName string
	Role     vsec.AuthLevel
}

func getKey() []byte {
	return []byte("valrrwwssffgsdgfksdjfghsdlgnsda")
}

func getAccessLevel(path string) (access vsec.AuthLevel, err error) {
	if len(path) > (accessPos+2) && path[accessPos] == 'r' {
		switch path[accessPos+1] {
		case '0':
			access = vsec.Super
		case '1':
			access = vsec.Admin
		case '2':
			access = vsec.Normal
		case '3':
			access = vsec.Monitor
		}
		access = vsec.Public
		err = fmt.Errorf("Invalid authorized URL: %s", path)
	}
	return access, err
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		var userRole, access vsec.AuthLevel
		access, err = getAccessLevel(ctx.Path())
		if err != nil {
			err = &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "Invalid URL",
				Inner:   err,
			}
			vlog.LogError("Net", err)
		}
		var userInfo JWTUserInfo
		userInfo, err = RetrieveUserInfo(ctx)
		if err != nil {
			err = &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "Invalid JWT toke found, does not have user info",
				Inner:   err,
			}
			vlog.LogError("Net", err)
		}
		if access < userRole {
			err = &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "",
				Inner:   err,
			}
			return err
		}
		if err == nil {
			ctx.Set("userID", userInfo.UserID)
			ctx.Set("userName", userInfo.UserName)
			err = next(ctx)
		}
		return vlog.LogError("Net", err)
	}
}

func dummyAuthenticator(params map[string]interface{}) (
	user *vsec.User, err error) {
	user = nil
	err = errors.New("No valid authenticator found")
	return user, err
}

func dummyAuthorizer(userID string) (role vsec.AuthLevel, err error) {
	err = errors.New("No valid authorizer found")
	return role, err
}

func login(ctx echo.Context) (err error) {
	defer func() {
		vlog.LogError("Net:Sec:API", err)
	}()
	msg := "Login successful"
	status := http.StatusOK
	var data map[string]interface{}
	userID := ""
	name := "" //user name is used for auditing
	creds := make(map[string]string)
	err = ctx.Bind(&creds)
	if err == nil {
		var user *vsec.User
		userID = creds["userID"]
		name = userID
		user, err = DoLogin(userID, creds["password"])
		if err == nil {
			if user.State == vsec.Active {
				token := jwt.New(jwt.SigningMethodHS256)
				claims := token.Claims.(jwt.MapClaims)
				name = user.FirstName + " " + user.LastName
				claims["userID"] = user.ID
				claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
				claims["access"] = user.Auth
				claims["userName"] = name
				var signed string
				signed, err = token.SignedString(getKey())
				if err == nil {
					data = make(map[string]interface{})
					data["token"] = signed
					data["user"] = user
				} else {
					msg = "Failed to sign token"
					status = http.StatusInternalServerError
				}
			} else {
				data = make(map[string]interface{})
				data["state"] = user.State
				msg = "User is not active"
				status = http.StatusUnauthorized
				err = errors.New(msg)
			}
		} else {
			msg = "Login failed"
			status = http.StatusUnauthorized
		}
	} else {
		msg = "Failed to read credentials from request"
		status = http.StatusBadRequest
	}
	//SHA1 encoded to avoid storing email in db
	ctx.Set("userID", vcmn.Hash(userID))
	ctx.Set("userName", name)
	AuditedSend(ctx, &Result{
		Status: status,
		Op:     "login",
		Msg:    msg,
		OK:     err == nil,
		Data:   data,
		Err:    err,
	})
	return vlog.LogError("Net:Sec:API", err)
}

//DoLogin - performs login using username and password
func DoLogin(userID string, password string) (*vsec.User, error) {
	//Check for password expiry and stuff
	params := make(map[string]interface{})
	params["userID"] = userID
	params["password"] = password
	user, err := authenticator(params)
	if err == nil && authorizer != nil {
		user.Auth, err = authorizer(user.ID)
	}
	return user, err
}

//RetrieveUserInfo - retrieves user information from JWT token
func RetrieveUserInfo(ctx echo.Context) (uinfo JWTUserInfo, err error) {
	success := false
	itk := ctx.Get("token")
	// vcmn.DumpJSON(itk)
	if tkn, ok := itk.(*jwt.Token); ok {
		if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
			var ok1, ok3 bool
			uinfo.UserID, ok1 = claims["userID"].(string)
			access, ok2 := claims["access"].(float64)
			uinfo.UserName, ok3 = claims["userName"].(string)
			uinfo.Role = vsec.AuthLevel(access)
			success = ok1 && ok2
			if !ok1 {
				vlog.Error("Net:Sec:API", "Invalid user ID in JWT")
			}
			if !ok2 {
				vlog.Error("Net:Sec:API", "Invalid access level in JWT")
			}
			if !ok3 {
				vlog.Error("Net:Sec:API", "Invalid user name in JWT")
			}
		}
	}
	if !success {
		err = errors.New("Could not find relevent information in JWT token")
	}
	return uinfo, err
}

//IsAdmin - returns true if user associated with request is an admin
func IsAdmin(ctx echo.Context) (yes bool) {
	yes = false
	uinfo, err := RetrieveUserInfo(ctx)
	if err == nil {
		yes = uinfo.Role <= vsec.Admin
	}
	return yes
}

//IsSuperUser - returns true if user is a super user
func IsSuperUser(ctx echo.Context) (yes bool) {
	yes = false
	uinfo, err := RetrieveUserInfo(ctx)
	vcmn.DumpJSON(uinfo)
	if err == nil {
		yes = uinfo.Role == vsec.Super
	}
	return yes
}

//IsNormalUser - returns true if user is a normal user
func IsNormalUser(ctx echo.Context) (yes bool) {
	yes = false
	uinfo, err := RetrieveUserInfo(ctx)
	if err == nil {
		yes = uinfo.Role <= vsec.Normal
	}
	return yes
}
