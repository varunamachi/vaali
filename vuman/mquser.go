package vuman

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
	// _, err = conn.C("user").Upsert(bson.M{"id": user.ID}, user)
	err = conn.C("user").Insert(user)
	return vlog.LogError("UMan:Mongo", err)
}

//UpdateUser - updates user in database
func UpdateUser(user *vsec.User) (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("user").Update(bson.M{"id": user.ID}, user)
	return vlog.LogError("UMan:Mongo", err)
}

//DeleteUser - deletes user with given user ID
func DeleteUser(userID string) (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("user").Remove(bson.M{"id": userID})
	return vlog.LogError("UMan:Mongo", err)
}

//GetUser - gets details of the user corresponding to ID
func GetUser(userID string) (user *vsec.User, err error) {
	conn := vdb.DefaultMongoConn()
	user = &vsec.User{}
	defer conn.Close()
	err = conn.C("user").Find(bson.M{"id": userID}).One(user)
	return user, vlog.LogError("UMan:Mongo", err)
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
	return users, vlog.LogError("UMan:Mongo", err)
}

//ResetPassword - sets password of a unauthenticated user
func ResetPassword(userID, oldPwd, newPwd string) (err error) {
	conn := vdb.DefaultMongoConn()
	defer func() {
		conn.Close()
		vlog.LogError("UMan:Mongo", err)
	}()
	if err != nil {
		return err
	}
	var newHash string
	newHash, err = passlib.Hash(newPwd)
	if err != nil {
		return err
	}
	if err = ValidateUser(userID, oldPwd); err != nil {
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
	var newHash string
	newHash, err = passlib.Hash(newPwd)
	if err == nil {
		conn := vdb.DefaultMongoConn()
		defer conn.Close()
		setPasswordHash(conn, userID, newHash)
	}
	return vlog.LogError("UMan:Mongo", err)
}

func setPasswordHash(conn *vdb.MongoConn, userID, hash string) (
	err error) {
	_, err = conn.C("secret").Upsert(
		bson.M{
			"userID": userID,
		},
		bson.M{
			"userID": userID,
			"phash":  hash,
		})
	return err
}

//ValidateUser - validates user ID and password
func ValidateUser(userID, password string) (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	if err == nil {
		secret := bson.M{}
		err = conn.C("secret").
			Find(bson.M{"userID": userID}).
			Select(bson.M{"phash": 1, "_id": 0}).
			One(&secret)
		if err == nil {
			storedHash, ok := secret["phash"].(string)
			if ok {
				var newHash string
				newHash, err = passlib.Verify(password, storedHash)
				if err == nil && newHash != "" {
					err = setPasswordHash(conn, userID, newHash)
				}
			} else {
				err = errors.New("Failed to varify password")
			}
		}
	}
	return vlog.LogError("UMan:Mongo", err)
}

//GetUserAuthLevel - gets user authorization level
func GetUserAuthLevel(userID string) (level vsec.AuthLevel, err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("user").
		Find(bson.M{"userID": userID}).
		Select(bson.M{"auth": 1}).
		One(&level)
	return level, vlog.LogError("UMan:Mongo", err)
}

//CreateFirstSuperUser - creates the first super user for the application
func CreateFirstSuperUser(user *vsec.User, password string) (err error) {
	defer func() {
		vlog.LogError("UMan:Mongo", err)
	}()
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	var count int
	count, _ = conn.C("user").Find(bson.M{"auth": 0}).Count()
	if count > 5 {
		err = errors.New("A super admin already exists, operation aborted")
		return err
	}
	err = CreateUser(user)
	if err != nil {
		return err
	}
	err = SetPassword(user.ID, password)
	return err
}