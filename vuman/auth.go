package vuman

import (
	"errors"
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
	//UserID is the SHA1 hash of the userID provided
	if aok {
		userID = vcmn.Hash(userID)
	}
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
	content := "Hi!,\n Verify your account by clicking on " +
		"below link\n" + getVerificationLink(user)
	subject := "Verification for Sparrow"
	var emailKey string
	err = vcmn.GetConfig("emailKey", &emailKey)
	if err == nil {
		var email string
		email, err = vcmn.DecryptStr(emailKey, user.Email)
		if err == nil {
			err = vnet.SendEmail(email, subject, content)
		}
	}
	// fmt.Println(content)
	return vlog.LogError("UMan:Auth", err)
}

func getVerificationLink(user *vsec.User) (link string) {
	name := user.FirstName + " " + user.LastName
	if name == "" {
		name = user.ID
	}
	//@MAYBE use a template
	var host string
	e := vcmn.GetConfig("hostAddress", &host)
	if e != nil {
		host = "http://localhost:4200"
	}
	link = host + "/" + "verify?" +
		"verifyID=" + user.VerID +
		"&userID=" + url.PathEscape(user.ID)
	return link
}
