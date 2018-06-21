package vmgo

import (
	"github.com/varunamachi/vaali/vcmn"
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

//GetAll - gets all the items from collection 'dtype' selected by filter & paged
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
	selector := GenerateSelector(filter)
	count, err = conn.C(dtype).
		Find(selector).
		Count()
	return count, vlog.LogError("DB:Mongo", err)
}

//GetAllWithCount - gets all the items from collection 'dtype' selected by
//filter & paged also gives the total count of items selected by filter
func GetAllWithCount(dtype string,
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

//FillFilterValues - Fills given filter descriptors with possible values when
//possible for a data type
// func FillFilterValues(dtype string, fds []*FilterSpec) (
// 	out []*FilterSpec) {
// 	conn := DefaultMongoConn()
// 	defer conn.Close()
// 	for _, fdesc := range fds {
// 		switch fdesc.Type {
// 		case Value:
// 			fallthrough
// 		case Array:
// 			sdata := make([]string, 0, 100)
// 			e := conn.C(dtype).Find(nil).Distinct(fdesc.Name, sdata)
// 			fdesc.Data = sdata
// 			vlog.LogError("DB:Mongo", e)
// 		case Date:
// 			var dr DateRange
// 			e := conn.C(dtype).Pipe([]bson.M{
// 				bson.M{
// 					"$group": bson.M{
// 						"_id": bson.M{},
// 						"from": bson.M{
// 							"$min": fdesc.Name,
// 						},
// 						"to": bson.M{
// 							"$max": fdesc.Name,
// 						},
// 					},
// 				},
// 			}).One(&dr)
// 			fdesc.Data = &dr
// 			vlog.LogError("DB:Mongo", e)
// 		}
// 	}
// 	out = fds
// 	return out
// }

//GenerateSelector - creates mongodb query for a generic filter
func GenerateSelector(filter *Filter) (selector bson.M) {
	queries := make([]bson.M, 0, 100)
	for key, values := range filter.Props {
		if len(values) == 1 {
			queries = append(queries, bson.M{key: values[0]})
		} else if len(values) > 1 {
			orProps := make([]bson.M, 0, len(values))
			for _, val := range values {
				orProps = append(orProps, bson.M{key: val})
			}
			queries = append(queries, bson.M{"$or": orProps})
		}
	}
	for field, val := range filter.Bools {
		if val != nil {
			queries = append(queries, bson.M{field: val})
		}
	}
	for field, dateRange := range filter.Dates {
		if dateRange.IsValid() {
			queries = append(queries,
				bson.M{
					field: bson.M{
						"$gte": dateRange.From,
						"$lte": dateRange.To,
					},
				},
			)
		}
	}
	for field, matcher := range filter.Lists {
		if len(matcher.Fields) != 0 {
			mode := "$in"
			if matcher.MatchAll {
				mode = "$all"
			}
			queries = append(queries, bson.M{
				field: bson.M{
					mode: matcher.Fields,
				},
			})
		}
	}
	if len(queries) != 0 {
		selector = bson.M{
			"$and": queries,
		}
	}
	return selector
}

//GetFilterValues - provides values associated the fields defined in filter spec
func GetFilterValues(
	dtype string,
	specs FilterSpecList) (values bson.M, err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	values = bson.M{}
	for _, spec := range specs {
		switch spec.Type {
		case Prop:
			fallthrough
		case Array:
			props := make([]string, 0, 100)
			err = conn.C(dtype).Find(nil).Distinct(spec.Field, &props)
			values[spec.Field] = props
		case Date:
			var drange vcmn.DateRange
			err = conn.C(dtype).Pipe([]bson.M{
				bson.M{
					"$group": bson.M{
						"_id": nil,
						"from": bson.M{
							"$max": spec.Field,
						},
						"to": bson.M{
							"$min": spec.Field,
						},
					},
				},
			}).One(&drange)
			values[spec.Field] = drange
		case Boolean:
		case Search:
		case Static:
		}
	}
	return values, vlog.LogError("DB:Mongo", err)
}
