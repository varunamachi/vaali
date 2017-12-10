package vsec

import (
	"github.com/varunamachi/vaali/vlog"
)

//CreateUser - creates user in database
func CreateUser(user *User) (err error) {

	return vlog.LogError("Sec:Mongo", err)
}

//UpdateUser - updates user in database
func UpdateUser(user *User) (err error) {
	return vlog.LogError("Sec:Mongo", err)
}

//DeleteUser - deletes user with given user ID
func DeleteUser(userID string) (err error) {
	return vlog.LogError("Sec:Mongo", err)
}

//GetUser - gets details of the user corresponding to ID
func GetUser(userID string) (err error) {
	return vlog.LogError("Sec:Mongo", err)
}

//GetAllUsers - gets all users based on offset and limit
func GetAllUsers(offset, limit int) (err error) {
	return vlog.LogError("Sec:Mongo", err)
}

//ResetPassword - sets password of a unauthenticated user
func ResetPassword(userID, oldPwd, newPwd string) (err error) {
	return vlog.LogError("Sec:Mongo", err)
}

//SetPassword - sets password of a already authenticated user, old password
//is not required
func SetPassword(userID, newPwd string) (err error) {
	return vlog.LogError("Sec:Mongo", err)
}
