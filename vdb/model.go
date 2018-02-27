package vdb

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/varunamachi/vaali/vcmn"
)

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

//FilterType - Type of filter item
type FilterType string

//Value - filter for a value
const Value FilterType = "value"

//Array - filter for an array
const Array FilterType = "array'"

//Date - filter for data range
const Date FilterType = "dateRange"

//Boolean - filter for boolean field
const Boolean FilterType = "boolean"

//FilterDesc - possible values for filters
type FilterDesc struct {
	Field string      `json:"field" bson:"field"`
	Name  string      `json:"field" bson:"field"`
	Type  FilterType  `json:"field" bson:"field"`
	Data  interface{} `json:"data" bson:"data"`
}

//DateRange - represents a date range
type DateRange struct {
	From time.Time `json:"from" bson:"from"`
	To   time.Time `json:"to" bson:"to"`
}

//StoredItem - represents a value that is stored in database and is
//compatible with generic queries and handlers. Any struct with a need to
//support generic CRUD operations must implement and register a factory
//method to return it
type StoredItem interface {
	ID() bson.ObjectId
	SetCreationInfo(at time.Time, by string)
	SetModInfo(at time.Time, by string)
}

//FactoryFunc - Function for creating an instance of data type
type FactoryFunc func() StoredItem

var mongoStore *store
var defaultDB = "vaali"
var factories = make(map[string]FactoryFunc)
