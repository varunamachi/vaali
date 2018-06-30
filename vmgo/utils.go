package vmgo

import (
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
