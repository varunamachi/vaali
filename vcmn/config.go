package vcmn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/mitchellh/mapstructure"

	"github.com/varunamachi/vaali/vlog"
)

var config = make(map[string]interface{})

func readConfig(dirPath, appName string) (err error) {
	path := dirPath + "/" + appName + ".conf.json"
	if ExistsAsFile(path) {
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

//LoadConfig - loads configuration for app with given appName. Searches for
//configuration file in standard locations and loads based all of them. If
//same config values is present in different files, value for the file that is
//loaded last is kept
func LoadConfig(appName string) {
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
	readConfig(GetExecDir(), appName)
}

//PrintConfig - prints the configuration
func PrintConfig() {
	for k, v := range config {
		fmt.Printf("%s: %v\n", k, v)
	}
}

//GetStringConfig - gets a value associated with config key
func GetStringConfig(key string) (value string) {
	value = config[key].(string)
	return value
}

//GetConfig - retrieves config value for the given key and populates the
//value argument given. If the key does not exist in the config map or
//if its not possible to populate value arg from retrieved value and error is
//returned
func GetConfig(key string, value interface{}) (err error) {
	if val, ok := config[key]; ok {
		if im, ok := val.(map[string]interface{}); ok {
			err = mapstructure.Decode(im, value)
		} else {
			err = fmt.Errorf("Config for key %s is not in expected format",
				key)
		}

	} else {
		err = fmt.Errorf("Config with key %s not found", key)
	}
	return err
}

//HasConfig - checks if a value exists in config for a key
func HasConfig(key string) (yes bool) {
	_, yes = config[key]
	return yes
}
