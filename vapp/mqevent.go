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
func GetEvents(offset, limit int, filter vdb.Filter) (
	total int, events []*vlog.Event, err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	events = make([]*vlog.Event, 0, limit)
	var selector bson.M
	selector, err = vdb.GenerateSelector(filter)
	if err == nil {
		q := conn.C("events").Find(selector).Sort("-time")
		total, err = q.Count()
		if err == nil {
			err = q.Skip(offset).Limit(limit).All(&events)
		}
	}
	return total, events, vlog.LogError("App:Event", err)
}

//GetEventFilterModel - gives filter event model generated from database
func GetEventFilterModel() (efm EventFilterModel, err error) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	efm = NewEventFilterModel()
	err = conn.C("users").Find(nil).Distinct("userID", &efm.UserIDs)
	if err == nil {
		err = conn.C("events").Find(nil).Distinct("op", &efm.EventTypes)
	}
	return efm, vlog.LogError("App:Event", err)
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
