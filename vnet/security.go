package vnet

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/varunamachi/vali/vsec"
)

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
		return err
	}
}
