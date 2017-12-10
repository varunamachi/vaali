package vnet

import (
	"errors"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vsec"
)

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

//RetrieveUserInfo - retrieves user information from JWT token
func RetrieveUserInfo(ctx echo.Context) (
	user string,
	role vsec.AuthLevel,
	err error) {
	success := false
	if tkn, ok := ctx.Get("token").(*jwt.Token); ok {
		if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
			var aok bool
			user, aok = claims["user"].(string)
			access, bok := claims["access"].(float64)
			role = vsec.AuthLevel(access)
			success = aok && bok
		}
	}
	if !success {
		err = errors.New("Could not find relevent information in JWT token")
	}
	return user, role, err
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		var userRole, access vsec.AuthLevel
		var user string
		access, err = getAccessLevel(ctx.Path())
		if err != nil {
			err = &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "Invalid URL",
				Inner:   err,
			}
			return err
		}
		user, userRole, err = RetrieveUserInfo(ctx)
		if err != nil {
			err = &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "Invalid JWT toke found, does not have user info",
				Inner:   err,
			}
			return err
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
			ctx.Set("userID", user)
			err = next(ctx)
		}
		return vlog.LogError("Net", err)
	}
}
