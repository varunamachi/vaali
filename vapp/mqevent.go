package vapp

import (
	"github.com/varunamachi/vaali/vdb"
	"github.com/varunamachi/vaali/vlog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoAuditor - stores event logs into database
func MongoAuditor(event *vlog.Event) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	conn.C("events").Insert(event)
}

//GetEvents - retrieves event entries based on filters
func GetEvents(offset, limit int, filter bson.M) (
	events []*vlog.Event, err error) {
	//@TODO define filter structure and the translate to query
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	events = make([]*vlog.Event, 0, limit)
	err = conn.C("events").
		// Find(filter).
		Find(nil).
		Sort("-time").
		Skip(offset).
		Limit(limit).
		All(&events)
	return events, vlog.LogError("App:Event", err)
}

//CreateIndices - creates mongoDB indeces for tables used for event logs
func CreateIndices() (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	err = conn.C("events").EnsureIndex(mgo.Index{
		Key:        []string{"time"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	return err
}

//CleanData - cleans event related data from database
func CleanData() (err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	_, err = conn.C("events").RemoveAll(bson.M{})
	return err
}
