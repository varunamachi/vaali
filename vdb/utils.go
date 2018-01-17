package vdb

import "gopkg.in/mgo.v2/bson"

//GenerateSelector - creates mongodb query for a generic filter
func GenerateSelector(filter Filter) (selector bson.M, err error) {
	selector = bson.M{}
	queries := make([]bson.M, 0, 100)
	for key, value := range filter.Fields {
		queries = append(queries, bson.M{key: value})
	}
	for _, date := range filter.Dates {
		queries = append(queries, bson.M{
			"$gte": date.From,
			"$lte": date.To,
		})
	}
	for _, matcher := range filter.Lists {
		mode := "$in"
		if matcher.MatchAll {
			mode = "$all"
		}
		queries = append(queries, bson.M{
			matcher.Name: bson.M{
				mode: matcher.Tags,
			},
		})
	}
	return selector, err
}
