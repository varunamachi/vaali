package vcmn

import (
	"fmt"
	"time"
)

//Version - represents version of the application
type Version struct {
	Major int `json:"major" bson:"major"`
	Minor int `json:"minor" bson:"minor"`
	Patch int `json:"patch" bson:"patch"`
}

//String - version to string
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

//DateRange - represents date ranges
type DateRange struct {
	// Name string    `json:"name" bson:"name"`
	From time.Time `json:"from" bson:"from"`
	To   time.Time `json:"to" bson:"to"`
}

//IsValid - returns true if both From and To dates are non-zero
func (r *DateRange) IsValid() bool {
	return !(r.From.IsZero() || r.To.IsZero())
}

//ParamType - type of the parameter
type ParamType int

const (
	//Boolean - bool parameter
	Boolean ParamType = iota

	//NumberRange - number range parameter
	NumberRange

	//Choice - parameter with choices
	Choice

	//Text - arbitrary string
	Text
)

//Pair - association of key and value
type Pair struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

//Range - integer range
type Range struct {
	Min int `json:"min" bson:"min"`
	Max int `json:"max" bson:"max"`
}

//Param - represents generic parameter
type Param struct {
	Name    string      `json:"name" bson:"name"`
	Type    ParamType   `json:"type" bson:"type"`
	Desc    string      `json:"desc" bson:"desc"`
	Value   interface{} `json:"value" bson:"value"`
	Range   Range       `json:"range" bson:"range"`
	Choices []Pair      `json:"choices" bson:"choices"`
	Default interface{} `json:"default" bson:"default"`
}
