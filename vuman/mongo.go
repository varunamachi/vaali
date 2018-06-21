package vuman

import (
	"errors"
	"time"

	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vmgo"
	"github.com/varunamachi/vaali/vsec"
	passlib "gopkg.in/hlandau/passlib.v1"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoStorage - mongodb storage for user information
type MongoStorage struct{}

//CreateUser - creates user in database
func (m *MongoStorage) CreateUser(user *vsec.User) (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	// _, err = conn.C("users").Upsert(bson.M{"id": user.ID}, user)
	err = conn.C("users").Insert(user)
	return vlog.LogError("UMan:Mongo", err)
}

//UpdateUser - updates user in database
func (m *MongoStorage) UpdateUser(user *vsec.User) (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("users").Update(bson.M{"id": user.ID}, user)
	return vlog.LogError("UMan:Mongo", err)
}

//DeleteUser - deletes user with given user ID
func (m *MongoStorage) DeleteUser(userID string) (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("users").Remove(bson.M{"id": userID})
	return vlog.LogError("UMan:Mongo", err)
}

//GetUser - gets details of the user corresponding to ID
func (m *MongoStorage) GetUser(userID string) (user *vsec.User, err error) {
	conn := vmgo.DefaultMongoConn()
	user = &vsec.User{}
	defer conn.Close()
	err = conn.C("users").Find(bson.M{"id": userID}).One(user)
	return user, vlog.LogError("UMan:Mongo", err)
}

//GetUsers - gets all users based on offset, limit and filter
func (m *MongoStorage) GetUsers(offset, limit int, filter *vmgo.Filter) (
	users []*vsec.User, err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	selector := vmgo.GenerateSelector(filter)
	users = make([]*vsec.User, 0, limit)
	err = conn.C("users").
		Find(selector).
		Sort("-created").
		Skip(offset).
		Limit(limit).
		All(&users)
	return users, vlog.LogError("UMan:Mongo", err)
}

//GetCount - gives the number of user selected by given filter
func (m *MongoStorage) GetCount(filter *vmgo.Filter) (count int, err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	selector := vmgo.GenerateSelector(filter)
	count, err = conn.C("users").Find(selector).Count()
	return users, vlog.LogError("UMan:Mongo", err)
}

//GetAllUsersWithCount - gets all users based on offset and limit, total count
//is also returned
// func (m *MongoStorage) GetAllUsersWithCount(offset, limit int) (
// 	total int, users []*vsec.User, err error) {
// 	conn := vmgo.DefaultMongoConn()
// 	defer conn.Close()
// 	users = make([]*vsec.User, 0, limit)
// 	q := conn.C("users").Find(bson.M{}).Sort("-created")
// 	total, err = q.Count()
// 	if err == nil {
// 		err = q.Skip(offset).Limit(limit).All(&users)
// 	}
// 	return total, users, vlog.LogError("UMan:Mongo", err)
// }

//GetUsersWithCount - Get users with total count
func (m *MongoStorage) GetUsersWithCount(
	offset, limit int, filter *vmgo.Filter) (
	total int, users []*vsec.User, err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	var selector bson.M
	selector = vmgo.GenerateSelector(filter)
	users = make([]*vsec.User, 0, limit)
	q := conn.C("users").Find(selector).Sort("-created")
	total, err = q.Count()
	if err == nil {
		err = q.Skip(offset).Limit(limit).All(&users)
	}
	return total, users, vlog.LogError("UMan:Mongo", err)
}

//ResetPassword - sets password of a unauthenticated user
func (m *MongoStorage) ResetPassword(
	userID, oldPwd, newPwd string) (err error) {
	conn := vmgo.DefaultMongoConn()
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
	if err = m.ValidateUser(userID, oldPwd); err != nil {
		err = errors.New("Could not match old password")
		return err
	}
	err = conn.C("secret").Update(
		bson.M{
			"userID": userID,
		},
		bson.M{
			"$set": bson.M{
				"phash": newHash,
			},
		},
	)
	return err
}

//SetPassword - sets password of a already authenticated user, old password
//is not required
func (m *MongoStorage) SetPassword(userID, newPwd string) (err error) {
	var newHash string
	newHash, err = passlib.Hash(newPwd)
	if err == nil {
		conn := vmgo.DefaultMongoConn()
		defer conn.Close()
		m.setPasswordHash(conn, userID, newHash)
	}
	return vlog.LogError("UMan:Mongo", err)
}

func (m *MongoStorage) setPasswordHash(conn *vmgo.MongoConn,
	userID, hash string) (err error) {
	_, err = conn.C("secret").Upsert(
		bson.M{
			"userID": userID,
		},
		bson.M{
			"$set": bson.M{
				"userID": userID,
				"phash":  hash,
			},
		})
	return err
}

//ValidateUser - validates user ID and password
func (m *MongoStorage) ValidateUser(userID, password string) (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
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
				err = m.setPasswordHash(conn, userID, newHash)
			}
		} else {
			err = errors.New("Failed to varify password")
		}
	}
	return vlog.LogError("UMan:Mongo", err)
}

//GetUserAuthLevel - gets user authorization level
func (m *MongoStorage) GetUserAuthLevel(
	userID string) (level vsec.AuthLevel, err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("users").
		Find(bson.M{"userID": userID}).
		Select(bson.M{"auth": 1}).
		One(&level)
	return level, vlog.LogError("UMan:Mongo", err)
}

//CreateSuperUser - creates the first super user for the application
func (m *MongoStorage) CreateSuperUser(
	user *vsec.User, password string) (err error) {
	defer func() {
		vlog.LogError("UMan:Mongo", err)
	}()
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	var count int
	count, _ = conn.C("users").Find(bson.M{"auth": 0}).Count()
	if count > 5 {
		err = errors.New("Super user limit exceeded")
		return err
	}
	err = updateUserInfo(user)
	if err != nil {
		return err
	}
	user.State = vsec.Active
	err = m.CreateUser(user)
	if err != nil {
		return err
	}
	err = m.SetPassword(user.ID, password)
	return err
}

//SetUserState - sets state of an user account
func (m *MongoStorage) SetUserState(
	userID string, state vsec.UserState) (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("users").Update(
		bson.M{
			"id": userID,
		},
		bson.M{
			"$set": bson.M{
				"state": state,
			},
		})
	return vlog.LogError("UMan:Mongo", err)
}

//VerifyUser - sets state of an user account to verified based on userID
//and verification ID
func (m *MongoStorage) VerifyUser(userID, verID string) (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("users").Update(
		bson.M{
			"$and": []bson.M{
				bson.M{"id": userID},
				bson.M{"verID": verID},
			},
		},
		bson.M{
			"$set": bson.M{
				"state":    vsec.Active,
				"verified": time.Now(),
			},
		})
	return vlog.LogError("UMan:Mongo", err)
}

//CreateIndices - creates mongoDB indeces for tables used for user management
func (m *MongoStorage) CreateIndices() (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("users").EnsureIndex(mgo.Index{
		Key:        []string{"id", "varfnID"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	})
	return err
}

//CleanData - cleans user management related data from database
func (m *MongoStorage) CleanData() (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	_, err = conn.C("users").RemoveAll(bson.M{})
	return err
}
