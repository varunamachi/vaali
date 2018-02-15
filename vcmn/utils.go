package vcmn

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/varunamachi/vaali/vlog"
)

//DumpJSON - dumps JSON representation of given data to stdout
func DumpJSON(o interface{}) {
	b, err := json.MarshalIndent(o, "", "    ")
	if err == nil {
		fmt.Println(string(b))
	} else {
		vlog.LogError("Cmn:Utils", err)
	}
}

//GetAsJSON - converts given data to JSON and returns as pretty printed
func GetAsJSON(o interface{}) (jstr string, err error) {
	b, err := json.MarshalIndent(o, "", "    ")
	if err == nil {
		jstr = string(b)
	}
	return jstr, vlog.LogError("Cmn:Utils", err)
}

//GetExecDir - gives absolute path of the directory in which the executable
//for the current application is present
func GetExecDir() (dirPath string) {
	execPath, err := os.Executable()
	if err == nil {
		dirPath = filepath.Dir(execPath)
	} else {
		vlog.LogError("Cmn:Utils", err)
	}

	return dirPath
}

//ExistsAsFile - checks if a regular file exists at given path. If a error
//occurs while stating whatever exists at given location, false is returned
func ExistsAsFile(path string) (yes bool) {
	stat, err := os.Stat(path)
	if err == nil && !stat.IsDir() {
		yes = true
	}
	return yes
}

//ExistsAsDir - checks if a directory exists at given path. If a error
//occurs while stating whatever exists at given location, false is returned
func ExistsAsDir(path string) (yes bool) {
	stat, err := os.Stat(path)
	if err == nil && stat.IsDir() {
		yes = true
	}
	return yes
}

//ErrString - returns the error string if the given error is not nil
func ErrString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

//FirstValid - returns the first error that is not nil
func FirstValid(errs ...error) (err error) {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
