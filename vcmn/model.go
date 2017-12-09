package vcmn

import "fmt"

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
