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
	filter *Filter,
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

//Count - counts the number of items for data type
func Count(dtype string, filter *Filter) (count int, err error) {
	//@TODO handle filters
	conn := DefaultMongoConn()
	defer conn.Close()
	var selector bson.M
	selector, err = GenerateSelector(filter)
	if err == nil {
		count, err = conn.C(dtype).
			Find(selector).
			Count()
	}
	return count, vlog.LogError("DB:Mongo", err)
}

//FillFilterValues - Fills given filter descriptors with possible values when
//possible for a data type
func FillFilterValues(dtype string, fds []*FilterDesc) (
	out []*FilterDesc) {
	conn := DefaultMongoConn()
	defer conn.Close()
	for _, fdesc := range fds {
		switch fdesc.Type {
		case Value:
			fallthrough
		case Array:
			sdata := make([]string, 0, 100)
			e := conn.C(dtype).Find(nil).Distinct(fdesc.Name, sdata)
			fdesc.Data = sdata
			vlog.LogError("DB:Mongo", e)
		case Date:
			var dr DateRange
			e := conn.C(dtype).Pipe([]bson.M{
				bson.M{
					"$group": bson.M{
						"_id": bson.M{},
						"from": bson.M{
							"$min": fdesc.Name,
						},
						"to": bson.M{
							"$max": fdesc.Name,
						},
					},
				},
			}).One(&dr)
			fdesc.Data = &dr
			vlog.LogError("DB:Mongo", e)
		}
	}
	out = fds
	return out
}
