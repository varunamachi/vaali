package vnet

import (
	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vsec"
)

//GetEndpoints - Export app security related APIs
func GetEndpoints() (endpoints []*Endpoint) {
	endpoints = []*Endpoint{
		&Endpoint{
			Method:   echo.POST,
			URL:      "login",
			Category: "security",
			Func:     login,
			Access:   vsec.Public,
		},
	}
	return endpoints
}
