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
			URL:      "uman/user",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     createUser,
			Comment:  "Create an user",
		},
		&vnet.Endpoint{
			Method:   echo.PUT,
			URL:      "uman/user",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     updateUser,
			Comment:  "Update an user",
		},
		&vnet.Endpoint{
			Method:   echo.DELETE,
			URL:      "uman/user/:userID",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     deleteUser,
			Comment:  "Delete an user",
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "uman/user/:userID",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     getUser,
			Comment:  "Get info about an user",
		},
		&vnet.Endpoint{
			Method:   echo.GET,
			URL:      "uman/user",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     getUsers,
			Comment:  "Get list of user & their details",
		},
		&vnet.Endpoint{
			Method:   echo.POST,
			URL:      "uman/user/password",
			Access:   vsec.Admin,
			Category: "user management",
			Func:     setPassword,
			Comment:  "Set password for an user",
		},
		&vnet.Endpoint{
			Method:   echo.PUT,
			URL:      "uman/user/password",
			Access:   vsec.Monitor,
			Category: "user management",
			Func:     resetPassword,
			Comment:  "Reset password",
		},
		&vnet.Endpoint{
			Method:   echo.POST,
			URL:      "/uman/user/self",
			Access:   vsec.Public,
			Category: "user management",
			Func:     registerUser,
			Comment:  "Registration for new user",
		},
		&vnet.Endpoint{
			Method:   echo.POST,
			URL:      "/uman/user/verify/:userID/:verID",
			Access:   vsec.Public,
			Category: "user management",
			Func:     verifyUser,
			Comment:  "Verify a registered account",
		},
		//@TODO implement BELOW - same as updateUser but can only update current
		//user
		// &vnet.Endpoint{
		// 	Method:   echo.PUT,
		// 	URL:      "/uman/user/self",
		// 	Access:   vsec.Public,
		// 	Category: "user management",
		// 	Func:     updateUserProfile,
		// },
	}
	return endpoints

}
