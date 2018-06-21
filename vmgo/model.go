package vmgo

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/varunamachi/vaali/vcmn"
)

//MongoConnOpts - options for connecting to a mongodb instance
type MongoConnOpts struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

//store - holds mongodb connection handle and information
type store struct {
	session *mgo.Session
	opts    []*MongoConnOpts
}

//CountList - paginated list returned from mongoDB along with total number of
//items in the list counted without pagination
type CountList struct {
	TotalCount int         `json:"total" bson:"total"`
	Data       interface{} `json:"data" bson:"data"`
}

//FilterType - Type of filter item
type FilterType string

//Prop - filter for a value
const Prop FilterType = "prop"

//Array - filter for an array
const Array FilterType = "array'"

//Date - filter for data range
const Date FilterType = "dateRange"

//Boolean - filter for boolean field
const Boolean FilterType = "boolean"

//Search - filter for search text field
const Search FilterType = "search"

//Constant - constant filter value
const Constant FilterType = "constant"

//Static - constant filter value
const Static FilterType = "static"

//FilterSpec - filter specification
type FilterSpec struct {
	Field string     `json:"field" bson:"field"`
	Name  string     `json:"name" bson:"name"`
	Type  FilterType `json:"type" bson:"type"`
}

//Matcher - matches the given fields. If MatchAll set to true all
//the elements of the fields array needs to be matched, otherwise only one element
//needs to match (minimum)
type Matcher struct {
	MatchAll bool     `json:"matchAll" bson:"matchAll"`
	Fields   []string `json:"fields" bson:"fields"`
}

//SearchField - contains search string and info for performing the search
// type SearchField struct {
// 	MatchAll  bool   `json:"matchAll" bson:"matchAll"`
// 	Regex     bool   `json:"regex" bson:"regex"`
// 	SearchStr string `json:"searchStr" bson:"searchStr"`
// }

//PropMatcher - matches props
type PropMatcher []interface{}

//Filter - generic filter used to filter data in any mongodb collection
type Filter struct {
	Props    map[string]PropMatcher    `json:"props" bson:"props"`
	Bools    map[string]interface{}    `json:"bools" bson:"bools"`
	Dates    map[string]vcmn.DateRange `json:"dates" bson:"dates"`
	Lists    map[string]Matcher        `json:"lists" bson:"lists"`
	Searches map[string]Matcher        `json:"searches" bson:"searches"`
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

//FilterSpecList - alias for array of filter specs
type FilterSpecList []*FilterSpec

//FactoryFunc - Function for creating an instance of data type
type FactoryFunc func() StoredItem

var mongoStore *store
var defaultDB = "vaali"
var factories = make(map[string]FactoryFunc)
