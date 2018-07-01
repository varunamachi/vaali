package vevt

import (
	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vmgo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoAuditor - handles application events and stores them in mongodb for audit
// purposes
type MongoAuditor struct{}

//LogEvent - stores event logs into database
func (m *MongoAuditor) LogEvent(event *Event) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	conn.C("events").Insert(event)
}

//GetEvents - retrieves event entries based on filters
func (m *MongoAuditor) GetEvents(offset, limit int, filter *vcmn.Filter) (
	total int, events []*Event, err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	events = make([]*Event, 0, limit)
	selector := vmgo.GenerateSelector(filter)
	q := conn.C("events").Find(selector).Sort("-time")
	total, err = q.Count()
	if err == nil {
		err = q.Skip(offset).Limit(limit).All(&events)
	}
	return total, events, vmgo.LogError("App:Event", err)
}

//CreateIndices - creates mongoDB indeces for tables used for event logs
func (m *MongoAuditor) CreateIndices() (err error) {
	conn := vmgo.DefaultMongoConn()
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
func (m *MongoAuditor) CleanData() (err error) {
	conn := vmgo.DefaultMongoConn()
	defer conn.Close()
	_, err = conn.C("events").RemoveAll(bson.M{})
	return err
}
