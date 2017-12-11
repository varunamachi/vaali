package vnet

import (
	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vsec"
	"gopkg.in/mgo.v2/bson"
)

//Endpoint - represents a REST endpoint with associated metadata
type Endpoint struct {
	OID      bson.ObjectId  `json:"_id"`
	Method   string         `json:"method"`
	URL      string         `json:"url"`
	Access   vsec.AuthLevel `json:"access"`
	Category string         `json:"cateogry"`
	Route    *echo.Route    `json:"route"`
	Func     echo.HandlerFunc
}

//Result - result of an API call
type Result struct {
	Op   string      `json:"op" bson:"op"`
	Msg  string      `json:"msg" bson:"msg"`
	OK   bool        `json:"ok" bson:"ok"`
	Data interface{} `json:"data" bson:"data"`
}
