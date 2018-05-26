package vsec

import (
	"time"

	"github.com/varunamachi/vaali/vmgo"
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

//Verfied - user account is verified by the user
var Verfied UserState = "verified"

//Active - user is active
var Active UserState = "active"

//Disabled - user account is disabled by an admin
var Disabled UserState = "disabled"

//Flagged - user account is flagged by a user
var Flagged UserState = "flagged"

//User - represents an user
type User struct {
	OID         bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	ID          string        `json:"id" bson:"id"`
	Email       string        `json:"email" bson:"email"`
	Auth        AuthLevel     `json:"auth" bson:"auth"`
	FirstName   string        `json:"firstName" bson:"firstName"`
	LastName    string        `json:"lastName" bson:"lastName"`
	State       UserState     `json:"state" bson:"state"`
	VerID       string        `json:"verID" bson:"verID"`
	PwdExpiry   time.Time     `json:"pwdExpiry" bson:"pwdExpiry"`
	Created     time.Time     `json:"created" bson:"created"`
	Modified    time.Time     `json:"modified" bson:"modified"`
	VerfiedDate time.Time     `json:"verified" bson:"verified"`
	Props       bson.M        `json:"props" bson:"props"`
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

//UserStorage - interface representing strategy to store and manage user
//information
type UserStorage interface {
	//CreateUser - creates user in database
	CreateUser(user *User) (err error)

	//UpdateUser - updates user in database
	UpdateUser(user *User) (err error)

	//DeleteUser - deletes user with given user ID
	DeleteUser(userID string) (err error)

	//GetUser - gets details of the user corresponding to ID
	GetUser(userID string) (user *User, err error)

	//GetAllUsers - gets all users based on offset and limit
	GetAllUsers(offset, limit int) (total int, users []*User, err error)

	//GetUsers - gives a list of users based on their state
	GetUsers(offset, limit int, filter *vmgo.Filter) (
		total int, users []*User, err error)

	//ResetPassword - sets password of a unauthenticated user
	ResetPassword(userID, oldPwd, newPwd string) (err error)

	//SetPassword - sets password of a already authenticated user, old password
	//is not required
	SetPassword(userID, newPwd string) (err error)

	//ValidateUser - validates user ID and password
	ValidateUser(userID, password string) (err error)

	//GetUserAuthLevel - gets user authorization level
	GetUserAuthLevel(userID string) (level AuthLevel, err error)

	//CreateSuperUser - creates the first super user for the application
	CreateSuperUser(user *User, password string) (err error)

	//SetUserState - sets state of an user account
	SetUserState(userID string, state UserState) (err error)

	//VerifyUser - sets state of an user account to verified based on userID
	//and verification ID
	VerifyUser(userID, verID string) (err error)

	//CreateIndices - creates mongoDB indeces for tables used for user management
	CreateIndices() (err error)

	//CleanData - cleans user management related data from database
	CleanData() (err error)
}
