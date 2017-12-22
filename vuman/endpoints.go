package vuman

import (
	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vnet"
	"github.com/varunamachi/vaali/vsec"
)

//GetEndpoints - gives REST endpoints for user manaagement
func GetEndpoints() (endpoints []*vnet.Endpoint) {
	endpoints = []*vnet.Endpoint{
		&vnet.Endpoint{
			Method:   echo.POST,
			URL:      "/uman/user",
			Access:   vsec.Public,
			Category: "user management",
			Func:     createUser,
		},
		&vnet.Endpoint{
			Method:   echo.PUT,
			URL:      "/uman/user",
			Access:   vsec.Public,
			Category: "user management",
			Func:     updateUser,
		},
		&vnet.Endpoint{
			Method:   echo.DELETE,
			URL:      "/uman/user/:userID",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     deleteUser,
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "/uman/user/:userID",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     getUser,
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "/uman/user",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     getUsers,
		},
		&vnet.Endpoint{
			Method:   echo.POST,
			URL:      "/uman/user/password",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     setPassword,
		},
		&vnet.Endpoint{
			Method:   echo.PUT,
			URL:      "/uman/user/password",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     resetPassword,
		},
	}
	return endpoints

}
