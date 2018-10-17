package vnet

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/varunamachi/vaali/vcmn"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vsec"
)

//GetJWTKey - gives a unique JWT key
func GetJWTKey() []byte {
	if len(jwtKey) == 0 {
		jwtKey, _ = uuid.NewV4().MarshalBinary()
	}
	return jwtKey
}

//Session - container for retrieving session & user information from JWT
type Session struct {
	UserID   string         `json:"userID"`
	UserName string         `json:"userName"`
	UserType string         `json:"userType"`
	Valid    bool           `json:"valid"`
	Role     vsec.AuthLevel `json:"role"`
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
			vlog.Error("Net", "URL: %s ERR: %v", ctx.Path(), err)
			err = &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "Invalid URL",
				Inner:   err,
			}
		}
		var userInfo Session
		userInfo, err = RetrieveSessionInfo(ctx)
		fmt.Println(err)
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
				claims["userType"] = "normal"
				var signed string
				key := GetJWTKey()
				signed, err = token.SignedString(key)
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
		Err:    vcmn.ErrString(err),
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

//GetToken - gets token from context or from header
func GetToken(ctx echo.Context) (token *jwt.Token, err error) {
	itk := ctx.Get("token")
	if itk != nil {
		var ok bool
		if token, ok = itk.(*jwt.Token); !ok {
			err = fmt.Errorf("Invalid token found in context")
		}
	} else {
		header := ctx.Request().Header.Get("Authorization")
		authSchemeLen := len("Bearer")
		if len(header) > authSchemeLen {
			tokStr := header[authSchemeLen+1:]
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				return GetJWTKey(), nil
			}
			token = new(jwt.Token)
			token, err = jwt.Parse(tokStr, keyFunc)
		} else {
			err = fmt.Errorf("Unexpected auth scheme used to JWT")
		}
	}
	return token, err
}

//RetrieveSessionInfo - retrieves session information from JWT token
func RetrieveSessionInfo(ctx echo.Context) (uinfo Session, err error) {
	success := true
	var token *jwt.Token
	if token, err = GetToken(ctx); err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			var access float64
			if uinfo.UserID, ok = claims["userID"].(string); !ok {
				vlog.Error("Net:Sec:API", "Invalid user ID in JWT")
				success = false
			}
			if uinfo.UserName, ok = claims["userName"].(string); !ok {
				vlog.Error("Net:Sec:API", "Invalid user name in JWT")
			}
			if uinfo.UserType, ok = claims["userType"].(string); !ok {
				vlog.Error("Net:Sec:API", "Invalid user type in JWT")
				success = false
			}
			if access, ok = claims["access"].(float64); !ok {
				vlog.Error("Net:Sec:API", "Invalid access level in JWT")
				success = false
			} else {
				uinfo.Role = vsec.AuthLevel(access)
			}
			uinfo.Valid = token.Valid
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
	uinfo, err := RetrieveSessionInfo(ctx)
	if err == nil {
		yes = uinfo.Role <= vsec.Admin
	}
	return yes
}

//IsSuperUser - returns true if user is a super user
func IsSuperUser(ctx echo.Context) (yes bool) {
	yes = false
	uinfo, err := RetrieveSessionInfo(ctx)
	vcmn.DumpJSON(uinfo)
	if err == nil {
		yes = uinfo.Role == vsec.Super
	}
	return yes
}

//IsNormalUser - returns true if user is a normal user
func IsNormalUser(ctx echo.Context) (yes bool) {
	yes = false
	uinfo, err := RetrieveSessionInfo(ctx)
	if err == nil {
		yes = uinfo.Role <= vsec.Normal
	}
	return yes
}
