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
