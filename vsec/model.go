package vsec

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//AuthLevel - authorization of an user
type AuthLevel int

const (
	//Super - super user
	Super AuthLevel = iota

	//Admin - application admin
	Admin

	//Normal - normal user
	Normal

	//Monitor - readonly user
	Monitor

	//Public - no authentication required
	Public
)

//User - represents an user
type User struct {
	OID       bson.ObjectId     `json:"_id" bson:"_id,omitempty"`
	ID        string            `json:"id" bson:"id"`
	Email     string            `json:"email" bson:"email"`
	Auth      AuthLevel         `json:"auth" bson:"auth"`
	FirstName string            `json:"firstName" bson:"firstName"`
	LastName  string            `json:"lastName" bson:"lastName"`
	Props     map[string]string `json:"props" bson:"props"`
	Created   time.Time         `json:"created" bson:"created"`
	Modified  time.Time         `json:"modified" bson:"modified"`
	PwdExpiry time.Time         `json:"pwdExpiry" bson:"pwdExpiry"`
}

//Group - group of users
type Group struct {
	OID   bson.ObjectId `json:"_id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Users []string      `json:"users" bson:"users"`
}

//String - string representation of AuthLevel
func (a AuthLevel) String() string {
	switch a {
	case Super:
		return "Super"
	case Admin:
		return "Admin"
	case Normal:
		return "Normal"
	case Monitor:
		return "Monitor"
	case Public:
		return "Public"

	}
	return "Unknown"
}
