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
			Access:   vsec.Monitor,
			Category: "generic",
			Func:     get,
			Comment:  "retrieve a resource of given type",
		},
		&Endpoint{
			Method:   echo.GET,
			URL:      "gen/:dataType/list",
			Access:   vsec.Monitor,
			Category: "generic",
			Func:     getAll,
			Comment:  "Retrieve a resource sub-list of given type",
		},
		&Endpoint{
			Method:   echo.GET,
			URL:      "gen/:dataType/count",
			Access:   vsec.Monitor,
			Category: "generic",
			Func:     count,
			Comment:  "Get count of items of data type",
		},
		&Endpoint{
			Method:   echo.GET,
			URL:      "gen/:dataType",
			Access:   vsec.Monitor,
			Category: "generic",
			Func:     getAllWithCount,
			Comment:  "Retrieve a resource sub-list of a type with total count",
		},
		&Endpoint{
			Method:   echo.GET,
			URL:      "gen/:dataType/fspec",
			Access:   vsec.Monitor,
			Category: "generic",
			Func:     getFilterValues,
			Comment:  "Get possible values for filter",
		},
		&Endpoint{
			Method:   echo.GET,
			URL:      "gen/:dataType/fvals",
			Access:   vsec.Monitor,
			Category: "generic",
			Func:     getFilterValuesX,
			Comment:  "Get possible values for filter",
		},
		// &Endpoint{
		// 	Method:   echo.GET,
		// 	URL:      "ping",
		// 	Access:   vsec.Public,
		// 	Category: "app",
		// 	Func:     ping,
		// 	Comment:  "Ping the server",
		// },
	}
	return endpoints
}
