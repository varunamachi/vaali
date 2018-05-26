package vuman

import "github.com/varunamachi/vaali/vsec"

//CleanData - cleans data from user storage
func CleanData() (err error) {
	return storage.CleanData()
}

//CreateIndices - creates required indices on the user data store
func CreateIndices() (err error) {
	return storage.CreateIndices()
}

//CreateSuperUser - creates super user using storage strategy
func CreateSuperUser(user *vsec.User, pwd string) (err error) {
	return storage.CreateSuperUser(user, pwd)
}

//SetPassword - sets user's password
func SetPassword(userID, pwd string) (err error) {
	return storage.SetPassword(userID, pwd)
}
