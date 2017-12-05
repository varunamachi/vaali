package vsec

import (
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

	//External - external user without authentication
	External
)

//User - represents an user
type User struct {
	OID       bson.ObjectId     `json:"_id" bson:"_id"`
	ID        string            `json:"id" bson:"id"`
	Email     string            `json:"email" bson:"email"`
	Auth      AuthLevel         `json:"auth" bson:"auth"`
	FirstName string            `json:"firstName" bson:"firstName"`
	LastName  string            `json:"lastName" bson:"lastName"`
	Props     map[string]string `json:"props" bson:"props"`
}

//Group - group of users
type Group struct {
	OID   bson.ObjectId `json:"_id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Users []string      `json:"users" bson:"users"`
}
