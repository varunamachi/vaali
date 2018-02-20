package vnet

import (
	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vsec"
	"gopkg.in/mgo.v2/bson"
)

//Authenticator - a function that is used to authenticate an user. The function
//takes map of parameters contents of which will differ based on actual function
//used
type Authenticator func(params map[string]interface{}) (*vsec.User, error)

//Authorizer - a function that will be used authorize an user
type Authorizer func(userID string) (vsec.AuthLevel, error)

//Endpoint - represents a REST endpoint with associated metadata
type Endpoint struct {
	OID      bson.ObjectId  `json:"_id"`
	Method   string         `json:"method"`
	URL      string         `json:"url"`
	Access   vsec.AuthLevel `json:"access"`
	Category string         `json:"cateogry"`
	Route    *echo.Route    `json:"route"`
	Comment  string         `json:"Comment"`
	Func     echo.HandlerFunc
}

//Result - result of an API call
type Result struct {
	Status int         `json:"status" bson:"status"`
	Op     string      `json:"op" bson:"op"`
	Msg    string      `json:"msg" bson:"msg"`
	OK     bool        `json:"ok" bson:"ok"`
	Err    error       `json:"error" bson:"error"`
	Data   interface{} `json:"data" bson:"data"`
}

//Options - options for initializing web APIs
type Options struct {
	RootName      string
	APIVersion    string
	Authenticator Authenticator
	Authorizer    Authorizer
}

//EmailConfig - configuration for sending email
type EmailConfig struct {
	AppEMail         string `json:"appEMail"`
	AppEMailPassword string `json:"appEMailPassword"`
	SMTPHost         string `json:"smtpHost"`
	SMTPPort         int    `json:"smtpPort"`
}

//M - map of string to interface
type M map[string]interface{}

//BinderFunc - a function that binds data struct to response body
type BinderFunc func(ctx echo.Context) (interface{}, error)

var categories = make(map[string][]*Endpoint)
var endpoints = make([]*Endpoint, 0, 200)
var e = echo.New()
var accessPos = 0
var rootPath = ""
var authenticator Authenticator
var authorizer Authorizer
