package vdb

import (
	"bytes"
	"strconv"

	"github.com/varunamachi/vaali/vlog"
	mgo "gopkg.in/mgo.v2"
)

var mongoStore *store
var defaultDB = "vaali"

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

//MongoConn - represents a mongdb connection
type MongoConn struct {
	*mgo.Database
}

//SetDefaultDB - sets the default DB
func SetDefaultDB(defDB string) {
	defaultDB = defDB
}

//Close - closes mongodb connection
func (conn *MongoConn) Close() {
	conn.Session.Close()
}

//toOptStr - converts mongodb options to a string that can be used as URL to
//connect to a mongodb instance
func toOptStr(options []*MongoConnOpts) string {
	var buf bytes.Buffer
	for i, opt := range options {
		//userName:password@host:port[,userName:password@host:port...]
		buf.WriteString("mongo://")
		if len(opt.User) != 0 {
			buf.WriteString(opt.User)
			buf.WriteString(":")
			buf.WriteString(opt.Password)
			buf.WriteString("@")
		}
		if len(opt.Host) != 0 {
			buf.WriteString(opt.Host)
		} else {
			buf.WriteString("localhost")
		}

		if opt.Port != 0 {
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(opt.Port))
		}
		if len(options) > 1 && i < len(options)-1 {
			//In case of multiple addresses
			buf.WriteString(",")
		}
	}
	return buf.String()
}

//ConnectSingle - connects to single instance of mongodb server
func ConnectSingle(opts *MongoConnOpts) (err error) {
	return Connect([]*MongoConnOpts{opts})
}

//Connect - connects to one or more mirrors of mongodb server
func Connect(opts []*MongoConnOpts) (err error) {
	var sess *mgo.Session
	optString := toOptStr(opts)
	sess, err = mgo.Dial(optString)
	if err == nil {
		sess.SetMode(mgo.Monotonic, true)
		mongoStore = &store{
			session: sess,
			opts:    opts,
		}
	}
	return vlog.LogError("DB:Mongo", err)
}

//NewMongoConn - creates a new connection to mogodb
func NewMongoConn(dbName string) (conn *MongoConn) {
	conn = &MongoConn{
		Database: mongoStore.session.Copy().DB(dbName),
	}
	return conn
}

//DefaultMongoConn - creates a connection to default DB
func DefaultMongoConn() *MongoConn {
	return NewMongoConn(defaultDB)
}

//CloseMongoConn - closes the mongodb connection
func CloseMongoConn() {
	mongoStore.session.Close()
}
