package vcmn

import (
	"fmt"
	"time"
)

//Version - represents version of the application
type Version struct {
	Major int `json:"major" bson:"major" sql:"major"`
	Minor int `json:"minor" bson:"minor" sql:"minor"`
	Patch int `json:"patch" bson:"patch" sql:"patch"`
}

//String - version to string
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

//DateRange - represents date ranges
type DateRange struct {
	// Name string    `json:"name" bson:"name"`
	From time.Time `json:"from" bson:"from" sql:"from"`
	To   time.Time `json:"to" bson:"to" sql:"to"`
}

//IsValid - returns true if both From and To dates are non-zero
func (r *DateRange) IsValid() bool {
	return !(r.From.IsZero() || r.To.IsZero())
}

//ParamType - type of the parameter
type ParamType int

const (
	//Bool - bool parameter
	Bool ParamType = iota

	//NumberRange - number range parameter
	NumberRange

	//Choice - parameter with choices
	Choice

	//Text - arbitrary string
	Text
)

//Pair - association of key and value
type Pair struct {
	Key   string `json:"key" bson:"key" sql:"key"`
	Value string `json:"value" bson:"value" sql:"value"`
}

//Range - integer range
type Range struct {
	Min int `json:"min" bson:"min" `
	Max int `json:"max" bson:"max"`
}

//Param - represents generic parameter
type Param struct {
	Name    string      `json:"name" bson:"name" sql:"name"`
	Type    ParamType   `json:"type" bson:"type" sql:"type"`
	Desc    string      `json:"desc" bson:"desc" sql:"desc"`
	Range   Range       `json:"range" bson:"range" sql:"range"`
	Choices []Pair      `json:"choices" bson:"choices" sql:"choices"`
	Default interface{} `json:"def" bson:"def" sql:"def"`
	// Value   interface{} `json:"value" bson:"value"`
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
	Props    map[string]PropMatcher `json:"props" bson:"props"`
	Bools    map[string]interface{} `json:"bools" bson:"bools"`
	Dates    map[string]DateRange   `json:"dates" bson:"dates"`
	Lists    map[string]Matcher     `json:"lists" bson:"lists"`
	Searches map[string]Matcher     `json:"searches" bson:"searches"`
}

//FilterSpecList - alias for array of filter specs
type FilterSpecList []*FilterSpec
