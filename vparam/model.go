package vparam

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

//Param - interface representing a parameter
type Param interface {
	Name() string
	Type() ParamType
	Description() string
	Value() interface{}
	SetValue(val interface{})
}

//BaseParam - basic parameter information
type BaseParam struct {
	Name string    `json:"name" bson:"name"`
	Type ParamType `json:"type" bson:"type"`
	Desc string    `json:"desc" bson:"desc"`
}

//BooleanParam - parameter with true/false value
type BooleanParam struct {
	Param
	Value   bool `json:"value" bson:"value"`
	Default bool `json:"value" bson:"value"`
}

//RangeParam - number range parameter
type RangeParam struct {
	Param
	Value   int  `json:"value" bson:"value"`
	Min     int  `json:"min" bson:"min"`
	Max     int  `json:"max" bson:"max"`
	Default bool `json:"value" bson:"value"`
}

//ChoiceParam - parameter with choices
type ChoiceParam struct {
	Param
	Value   string `json:"value" bson:"value"`
	Choices []Pair `json:"choices" bson:"choices"`
	Default int    `json:"default" bson:"default"`
}
