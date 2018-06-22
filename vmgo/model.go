package vmgo

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoConnOpts - options for connecting to a mongodb instance
type MongoConnOpts struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

//store - holds mongodb connection handle and information
type store struct {
	session *mgo.Session
	opts    []*MongoConnOpts
}

//StoredItem - represents a value that is stored in database and is
//compatible with generic queries and handlers. Any struct with a need to
//support generic CRUD operations must implement and register a factory
//method to return it
type StoredItem interface {
	ID() bson.ObjectId
	SetCreationInfo(at time.Time, by string)
	SetModInfo(at time.Time, by string)
}

//FactoryFunc - Function for creating an instance of data type
type FactoryFunc func() StoredItem

var mongoStore *store
var defaultDB = "vaali"
var factories = make(map[string]FactoryFunc)
