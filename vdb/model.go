package vdb

import "time"

//ArrayMatcher - matches elements of an array. If MatchAll set to true all
//the elements of the Tags array needs to be matched, otherwise only one element
//needs to match (minimum)
type ArrayMatcher struct {
	Name     string   `json:"name" bson:"name"`
	MatchAll bool     `json:"matchAll" bson:"matchAll"`
	Tags     []string `json:"tags" bson:"tags"`
}

//DateRange - represents date ranges
type DateRange struct {
	From time.Time `json:"from" bson:"from"`
	To   time.Time `json:"from" bson:"from"`
}

//Filter - generic filter used to filter data in any mongodb collection
type Filter struct {
	Fields map[string][]string `json:"fields" bson:"fields"`
	Dates  []DateRange         `json:"dates" bson:"dates"`
	Lists  []ArrayMatcher      `json:"lists" bson:"lists"`
}

//CountList - paginated list returned from mongoDB along with total number of
//items in the list counted without pagination
type CountList struct {
	TotalCount int         `json:"total" bson:"total"`
	Data       interface{} `json:"data" bson:"data"`
}
