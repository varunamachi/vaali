package vmgo

import (
	"runtime"

	"gopkg.in/mgo.v2"

	"github.com/varunamachi/vaali/vlog"
)

//RegisterFactory - registers factory functions for a data type
func RegisterFactory(dataType string, ff FactoryFunc) {
	factories[dataType] = ff
}

//Instance - creates and returns an instance of given data type
func Instance(dataType string) StoredItem {
	if ff, found := factories[dataType]; found {
		return ff()
	}
	vlog.Error("Generic:Inst", "Could not find factory for %s", dataType)
	return nil
}

//LogError - if error is not mog.ErrNotFound return null otherwise log the
//error and return the given error
func LogError(module string, err error) (out error) {
	if err != nil && err != mgo.ErrNotFound {
		_, file, line, _ := runtime.Caller(1)
		vlog.Error(module, "%s -- %s @ %d",
			err.Error(),
			file,
			line)
		out = err
	} else {
		err = nil
	}
	return out
}
