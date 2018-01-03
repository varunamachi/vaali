package vapp

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/varunamachi/vaali/vcmn"

	"github.com/varunamachi/vaali/vlog"
)

var config = make(map[string]string)

func readConfig(dirPath, appName string) (err error) {
	path := dirPath + "/" + appName + ".conf.json"
	if vcmn.ExistsAsFile(path) {
		var raw []byte
		raw, err = ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(raw, &config)
		}
		if err == nil {
			vlog.Info("App:Config", "Loaded config from %s", path)
		}
	}
	return err
}

func loadConfig(appName string) {
	switch runtime.GOOS {
	case "linux":
		readConfig("/etc/", appName)
		readConfig(os.ExpandEnv("$HOME"), appName)
	case "windows":
		readConfig(os.ExpandEnv("$ALLUSERSPROFILE"), appName)
		readConfig(os.ExpandEnv("$APPDATA"), appName)
	default:
		vlog.Warn("App:Config", "Unsupported operating system")
	}
	readConfig(vcmn.GetExecDir(), appName)
}

//GetConfig - gets a value associated with config key
func GetConfig(key string) (value string) {
	return config[key]
}

//GetConfigDef - gets a value associated with given key, if not found return
//default vlue
func GetConfigDef(key, def string) (value string) {
	var ok bool
	if value, ok = config[key]; ok {
		return config[key]
	}
	return def
}

//HasConfig - checks if a value exists in config for a key
func HasConfig(key string) (yes bool) {
	_, yes = config[key]
	return yes
}
