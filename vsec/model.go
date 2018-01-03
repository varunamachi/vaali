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

//UserState - state of the user account
type UserState string

//Active - user is active
var Active UserState = "active"

//Disabled - user account is disabled by an admin
var Disabled UserState = "disabled"

//Flagged - user account is flagged by a user
var Flagged UserState = "flagged"

//User - represents an user
type User struct {
	OID       bson.ObjectId     `json:"_id" bson:"_id,omitempty"`
	ID        string            `json:"id" bson:"id"`
	Email     string            `json:"email" bson:"email"`
	Auth      AuthLevel         `json:"auth" bson:"auth"`
	FirstName string            `json:"firstName" bson:"firstName"`
	LastName  string            `json:"lastName" bson:"lastName"`
	State     bool              `json:"state" bson:"state"`
	Created   time.Time         `json:"created" bson:"created"`
	Modified  time.Time         `json:"modified" bson:"modified"`
	PwdExpiry time.Time         `json:"pwdExpiry" bson:"pwdExpiry"`
	Props     map[string]string `json:"props" bson:"props"`
	VarfnID   string            `json:"varfnID" bson:"varfnID"`
}

//Group - group of users
type Group struct {
	OID   bson.ObjectId `json:"_id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Users []string      `json:"users" bson:"users"`
}

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
