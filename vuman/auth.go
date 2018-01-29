package vuman

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vnet"

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
	return user, vlog.LogError("UMan:Auth", err)
}

func sendVerificationMail(user *vsec.User) (err error) {

	name := user.FirstName + " " + user.LastName
	if name == "" {
		name = user.ID
	}
	//@MAYBE use a template
	var host string
	e := vcmn.GetConfig("emailHostAddress", &host)
	if e != nil {
		host = "localhost:80"
	}
	content := "Hello " + name + ",\n Verify your account by clicking on " +
		"below link\n" +
		"http://" +
		host + "/" +
		vnet.GetRootPath() +
		"/uman/user/verify/" +
		url.PathEscape(user.ID) +
		"/" +
		user.VerID
	// subject := "Verification for Sparrow"
	// err = vnet.SendEmail(user.Email, subject, content)
	fmt.Println(content)
	return vlog.LogError("UMan:Auth", err)
}
