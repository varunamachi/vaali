package vuman

import "github.com/varunamachi/vaali/vsec"

var storage vsec.UserStorage

//GetStorage - gets the underlying storage implementation
func GetStorage() vsec.UserStorage {
	return storage
}

//SetStorageStrategy - sets user storage strategy
func SetStorageStrategy(srg vsec.UserStorage) {
	storage = srg
}
