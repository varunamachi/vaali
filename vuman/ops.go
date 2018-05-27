package vuman

import "github.com/varunamachi/vaali/vsec"

//GetStorage - gets the underlying storage implementation
func GetStorage() vsec.UserStorage {
	return storage
}
