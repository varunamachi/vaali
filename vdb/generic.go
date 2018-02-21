package vdb

import (
	"github.com/varunamachi/vaali/vlog"
	"gopkg.in/mgo.v2/bson"
)

//Create - creates an record in 'dtype' collection
func Create(dtype string, value interface{}) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).Insert(value)
	return vlog.LogError("DB:Mongo", err)
}

//Update - updates the records in 'dtype' collection which are matched by
//the matcher query
func Update(dtype string, matcher bson.M, value interface{}) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).Update(matcher, value)
	return vlog.LogError("DB:Mongo", err)
}

//Delete - deletes record matched by the matcher from collection 'dtype'
func Delete(dtype string, matcher bson.M) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).Remove(matcher)
	return vlog.LogError("DB:Mongo", err)
}

//Get - gets a record matched by given matcher from collection 'dtype'
func Get(dtype string, matcher bson.M, out interface{}) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).Find(matcher).One(out)
	return vlog.LogError("DB:Mongo", err)
}

//GetAll - gets all the items from collection 'dtype'
func GetAll(dtype string,
	sortFiled string,
	offset int,
	limit int,
	out interface{}) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).
		Find(nil).
		Sort(sortFiled).
		Skip(offset).
		Limit(limit).
		All(out)
	return vlog.LogError("DB:Mongo", err)
}
