package vnet

import (
	"github.com/labstack/echo"
	"github.com/varunamachi/vali/vsec"
	"gopkg.in/mgo.v2/bson"
)

//EndpointFunc -
// type EndpointFunc func(ctx echo.Context) (err error)

//Endpoint -
type Endpoint struct {
	OID      bson.ObjectId  `json:"_id"`
	Method   string         `json:"method"`
	URL      string         `json:"url"`
	Access   vsec.AuthLevel `json:"access"`
	Category string         `json:"cateogry"`
	Route    *echo.Route    `json:"route"`
	Func     echo.HandlerFunc
}
