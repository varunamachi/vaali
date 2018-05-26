package vuman

import (
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vmgo"
	"github.com/varunamachi/vaali/vsec"
)

//PGStorage - postgresql storage for user information
type PGStorage struct{}

//CreateUser - creates user in database
func (p *PGStorage) CreateUser(user *vsec.User) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//UpdateUser - updates user in database
func (p *PGStorage) UpdateUser(user *vsec.User) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//DeleteUser - deletes user with given user ID
func (p *PGStorage) DeleteUser(userID string) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//GetUser - gets details of the user corresponding to ID
func (p *PGStorage) GetUser(userID string) (user *vsec.User, err error) {
	return user, vlog.LogError("UMan:PGSQL", err)
}

//GetAllUsers - gets all users based on offset and limit
func (p *PGStorage) GetAllUsers(offset, limit int) (
	total int, users []*vsec.User, err error) {
	return total, users, vlog.LogError("UMan:PGSQL", err)
}

//GetUsers - gives a list of users based on their state
func (p *PGStorage) GetUsers(
	offset, limit int, filter *vmgo.Filter) (
	total int, users []*vsec.User, err error) {
	return total, users, vlog.LogError("UMan:PGSQL", err)
}

//ResetPassword - sets password of a unauthenticated user
func (p *PGStorage) ResetPassword(
	userID, oldPwd, newPwd string) (err error) {
	return err
}

//SetPassword - sets password of a already authenticated user, old password
//is not required
func (p *PGStorage) SetPassword(userID, newPwd string) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

func (p *PGStorage) setPasswordHash(conn *vmgo.MongoConn,
	userID, hash string) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//ValidateUser - validates user ID and password
func (p *PGStorage) ValidateUser(userID, password string) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//GetUserAuthLevel - gets user authorization level
func (p *PGStorage) GetUserAuthLevel(
	userID string) (level vsec.AuthLevel, err error) {
	return level, vlog.LogError("UMan:PGSQL", err)
}

//CreateSuperUser - creates the first super user for the application
func (p *PGStorage) CreateSuperUser(
	user *vsec.User, password string) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//SetUserState - sets state of an user account
func (p *PGStorage) SetUserState(
	userID string, state vsec.UserState) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//VerifyUser - sets state of an user account to verified based on userID
//and verification ID
func (p *PGStorage) VerifyUser(userID, verID string) (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//CreateIndices - creates mongoDB indeces for tables used for user management
func (p *PGStorage) CreateIndices() (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}

//CleanData - cleans user management related data from database
func (p *PGStorage) CleanData() (err error) {
	return vlog.LogError("UMan:PGSQL", err)
}
