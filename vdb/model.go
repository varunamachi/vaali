package vdb

import "github.com/varunamachi/vaali/vcmn"

//ArrayMatcher - matches elements of an array. If MatchAll set to true all
//the elements of the Tags array needs to be matched, otherwise only one element
//needs to match (minimum)
type ArrayMatcher struct {
	// Name     string   `json:"name" bson:"name"`
	MatchAll bool     `json:"matchAll" bson:"matchAll"`
	Tags     []string `json:"tags" bson:"tags"`
}

//Filter - generic filter used to filter data in any mongodb collection
type Filter struct {
	Fields     map[string][]interface{}  `json:"fields" bson:"fields"`
	BoolFields map[string]bool           `json:"boolFields" bson:"boolFields"`
	Dates      map[string]vcmn.DateRange `json:"dates" bson:"dates"`
	Lists      map[string]ArrayMatcher   `json:"lists" bson:"lists"`
}

//CountList - paginated list returned from mongoDB along with total number of
//items in the list counted without pagination
type CountList struct {
	TotalCount int         `json:"total" bson:"total"`
	Data       interface{} `json:"data" bson:"data"`
}
