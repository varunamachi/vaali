package vuman

import (
	"errors"

	"github.com/varunamachi/vaali/vlog"

	"github.com/varunamachi/vaali/vsec"
)

func getUserIDPassword(params map[string]interface{}) (
	userID string, password string, err error) {
	var aok, bok bool
	userID, aok = params["userID"].(string)
	password, bok = params["password"].(string)
	if !aok || !bok {
		err = errors.New("Authorization, Invalid credentials provided")
	}
	return userID, password, err
}

//MongoAuthenticator - authenticator that uses user information stored in
//mongo DB to authenticate userID and password given
func MongoAuthenticator(params map[string]interface{}) (
	user *vsec.User, err error) {
	var userID, password string
	userID, password, err = getUserIDPassword(params)
	if err == nil {
		err = ValidateUser(userID, password)
		if err == nil {
			user, err = GetUser(userID)
		}
	}
	return user, vlog.LogError("UMan:Mongo:Auth", err)
}
