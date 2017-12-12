package vsec

import (
	"errors"

	"github.com/varunamachi/vaali/vdb"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vsec"
	passlib "gopkg.in/hlandau/passlib.v1"
	"gopkg.in/mgo.v2/bson"
)

//CreateUser - creates user in database
func CreateUser(user *vsec.User) (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	_, err = conn.C("user").Upsert(bson.M{"id": user.ID}, user)
	return vlog.LogError("Sec:Mongo", err)
}

//UpdateUser - updates user in database
func UpdateUser(user *vsec.User) (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("user").Update(bson.M{"id": user.ID}, user)
	return vlog.LogError("Sec:Mongo", err)
}

//DeleteUser - deletes user with given user ID
func DeleteUser(userID string) (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("user").Remove(bson.M{"id": userID})
	return vlog.LogError("Sec:Mongo", err)
}

//GetUser - gets details of the user corresponding to ID
func GetUser(userID string) (user *vsec.User, err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("user").Find(bson.M{"id": userID}).One(user)
	return user, vlog.LogError("Sec:Mongo", err)
}

//GetAllUsers - gets all users based on offset and limit
func GetAllUsers(offset, limit int) (users []*vsec.User, err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	users = make([]*vsec.User, 0, limit)
	err = conn.
		C("user").
		Find(bson.M{}).
		Sort("-created").
		Skip(offset).
		Limit(limit).
		All(&users)
	return users, vlog.LogError("Sec:Mongo", err)
}

//ResetPassword - sets password of a unauthenticated user
func ResetPassword(userID, oldPwd, newPwd string) (err error) {
	var oldHash, newHash string
	oldHash, err = passlib.Hash(oldPwd)
	conn := vdb.DefaultMongoConn()
	defer func() {
		conn.Close()
		vlog.LogError("Sec:Mongo", err)
	}()
	if err != nil {
		return err
	}
	newHash, err = passlib.Hash(newPwd)
	if err != nil {
		return err
	}
	storedPass := ""
	err = conn.C("secret").
		Find(bson.M{"userID": userID}).
		Select(bson.M{"phash": 1}).
		One(&storedPass)
	if err != nil || oldPwd == "" || oldHash != storedPass {
		err = errors.New("Could not match old password")
		return err
	}
	err = conn.C("secret").Update(
		bson.M{"userID": userID},
		bson.M{"phash": newHash},
	)
	return err
}

//SetPassword - sets password of a already authenticated user, old password
//is not required
func SetPassword(userID, newPwd string) (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	var newHash string
	newHash, err = passlib.Hash(newPwd)
	if err == nil {
		err = conn.C("secret").Update(
			bson.M{"userID": userID},
			bson.M{"phash": newHash},
		)
	}
	return vlog.LogError("Sec:Mongo", err)
}
