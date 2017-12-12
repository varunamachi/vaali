package vcmn

import "encoding/json"
import "fmt"
import "github.com/varunamachi/vaali/vlog"

//DumpJSON - dumps JSON representation of given data to stdout
func DumpJSON(o interface{}) {
	b, err := json.MarshalIndent(o, "", "    ")
	if err == nil {
		fmt.Println(string(b))
	} else {
		vlog.LogError("Cmd:JSON", err)
	}
}

//GetAsJSON - converts given data to JSON and returns as pretty printed
func GetAsJSON(o interface{}) (jstr string, err error) {
	b, err := json.MarshalIndent(o, "", "    ")
	if err == nil {
		jstr = string(b)
	}
	return jstr, vlog.LogError("Cmd:JSON", err)
}
