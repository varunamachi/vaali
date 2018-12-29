package vmgo

import (
	"github.com/varunamachi/vaali/vcmn"
	"gopkg.in/mgo.v2/bson"
)

//Create - creates an record in 'dtype' collection
func Create(dtype string, value interface{}) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).Insert(value)
	return LogError("DB:Mongo", err)
}

//Update - updates the records in 'dtype' collection which are matched by
//the matcher query
func Update(dtype string, matcher bson.M, value interface{}) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).Update(matcher, value)
	return LogError("DB:Mongo", err)
}

//Delete - deletes record matched by the matcher from collection 'dtype'
func Delete(dtype string, matcher bson.M) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).Remove(matcher)
	return LogError("DB:Mongo", err)
}

//Get - gets a record matched by given matcher from collection 'dtype'
func Get(dtype string, matcher bson.M, out interface{}) (err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).Find(matcher).One(out)
	return LogError("DB:Mongo", err)
}

//GetAll - gets all the items from collection 'dtype' selected by filter & paged
func GetAll(dtype string,
	sortFiled string,
	offset int,
	limit int,
	filter *vcmn.Filter,
	out interface{}) (err error) {
	selector := GenerateSelector(filter)
	conn := DefaultMongoConn()
	defer conn.Close()
	err = conn.C(dtype).
		Find(selector).
		Sort(sortFiled).
		Skip(offset).
		Limit(limit).
		All(out)
	return LogError("DB:Mongo", err)
}

//Count - counts the number of items for data type
func Count(dtype string, filter *vcmn.Filter) (count int, err error) {
	//@TODO handle filters
	conn := DefaultMongoConn()
	defer conn.Close()
	selector := GenerateSelector(filter)
	count, err = conn.C(dtype).
		Find(selector).
		Count()
	return count, LogError("DB:Mongo", err)
}

//GetAllWithCount - gets all the items from collection 'dtype' selected by
//filter & paged also gives the total count of items selected by filter
func GetAllWithCount(dtype string,
	sortFiled string,
	offset int,
	limit int,
	filter *vcmn.Filter,
	out interface{}) (count int, err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	selector := GenerateSelector(filter)
	q := conn.C(dtype).Find(selector)
	count, err = q.Count()
	if err == nil {
		err = q.Sort(sortFiled).
			Skip(offset).
			Limit(limit).
			All(out)
	}
	return count, LogError("DB:Mongo", err)
}

//GenerateSelector - creates mongodb query for a generic filter
func GenerateSelector(filter *vcmn.Filter) (selector bson.M) {
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
	// vcmn.DumpJSON(queries)
	return selector
}

//GetFilterValues - provides values associated the fields defined in filter spec
func GetFilterValues(
	dtype string,
	specs vcmn.FilterSpecList) (values bson.M, err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	values = bson.M{}
	for _, spec := range specs {
		switch spec.Type {
		case vcmn.Prop:
			fallthrough
		case vcmn.Array:
			props := make([]string, 0, 100)
			err = conn.C(dtype).Find(nil).Distinct(spec.Field, &props)
			values[spec.Field] = props
		case vcmn.Date:
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
		case vcmn.Boolean:
		case vcmn.Search:
		case vcmn.Static:
		}
	}
	return values, LogError("DB:Mongo", err)
}

//GetFilterValuesX - get values for filter based on given filter
func GetFilterValuesX(
	dtype string,
	field string,
	specs vcmn.FilterSpecList,
	filter *vcmn.Filter) (values bson.M, err error) {
	conn := DefaultMongoConn()
	defer conn.Close()
	facet := bson.M{}
	for _, spec := range specs {
		if spec.Field != field {
			switch spec.Type {
			case vcmn.Prop:
				facet[spec.Field] = []bson.M{
					bson.M{
						"$sortByCount": "$" + spec.Field,
					},
				}
			case vcmn.Array:
				fd := "$" + spec.Field
				facet[spec.Field] = []bson.M{
					bson.M{
						"$unwind": fd,
					},
					bson.M{
						"$sortByCount": fd,
					},
				}
			case vcmn.Date:
			case vcmn.Boolean:
			case vcmn.Search:
			case vcmn.Static:
			}
		}
	}
	var selector bson.M
	if filter != nil {
		selector = GenerateSelector(filter)
	}
	values = bson.M{}
	err = conn.C(dtype).Pipe([]bson.M{
		bson.M{
			"$match": selector,
		},
		bson.M{
			"$facet": facet,
		},
	}).One(&values)
	return values, LogError("DB:Mongo", err)
}
