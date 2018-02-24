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
			Comment:  "Login to application",
		},
		&Endpoint{
			Method:   echo.POST,
			URL:      "gen/:dataType",
			Access:   vsec.Normal,
			Category: "generic",
			Func:     create,
			Comment:  "Create a resource of given type",
		},
		&Endpoint{
			Method:   echo.PUT,
			URL:      "gen/:dataType",
			Access:   vsec.Normal,
			Category: "generic",
			Func:     update,
			Comment:  "Update a resource of given type",
		},
		&Endpoint{
			Method:   echo.DELETE,
			URL:      "gen/:dataType/:id",
			Access:   vsec.Normal,
			Category: "generic",
			Func:     delete,
			Comment:  "Delete a resource of given type",
		},
		&Endpoint{
			Method:   echo.GET,
			URL:      "gen/:dataType/:id",
			Access:   vsec.Normal,
			Category: "generic",
			Func:     get,
			Comment:  "retrieve a resource of given type",
		},
		&Endpoint{
			Method:   echo.GET,
			URL:      "gen/:dataType",
			Access:   vsec.Normal,
			Category: "generic",
			Func:     getAll,
			Comment:  "Retrieve a resource sub-list of given type",
		},
		&Endpoint{
			Method:   echo.GET,
			URL:      "gen/:dataType/count",
			Access:   vsec.Normal,
			Category: "generic",
			Func:     count,
			Comment:  "Create an resource of given type",
		},
	}
	return endpoints
}
