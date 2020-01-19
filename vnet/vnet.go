package vnet

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vsec"
)

//AddEndpoint - registers an REST endpoint
func AddEndpoint(ep *Endpoint) {
	endpoints = append(endpoints, ep)
}

//AddEndpoints - registers multiple REST categories
func AddEndpoints(eps ...*Endpoint) {
	for _, ep := range eps {
		AddEndpoint(ep)
	}
}

// ModifiedHTTPErrorHandler is the default HTTP error handler. It sends a
// JSON response with status code. [Modefied from echo.DefaultHTTPErrorHandler]
func ModifiedHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Inner != nil {
			msg = fmt.Sprintf("%v, %v", err, he.Inner)
		}
	} else if e.Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = echo.Map{"message": msg}
	}

	vlog.LogError("Net:HTTP", err)

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			vlog.LogError("Net:HTTP", err)
		}
	}
}

//InitWithOptions - initializes all the registered endpoints
func InitWithOptions(opts Options) {
	e.HideBanner = true
	e.HTTPErrorHandler = ModifiedHTTPErrorHandler
	e.Use(middleware.Recover())

	// TODO - enable based on env variable - better implement a custom
	// middleware to log only errors
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "[ACCSS] [Net:HTTP] ${status} : ${method} => ${path}\n",
	// }))
	//Add middleware
	authenticator = opts.Authenticator
	authorizer = opts.Authorizer

	//rootPath is a package variable
	rootPath = opts.RootName + "/api/v" + opts.APIVersion + "/"
	accessPos = len(rootPath) + len("in/")
	root := e.Group(rootPath)
	in := root.Group("in/")

	//For checking token
	in.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: GetJWTKey(),
		ContextKey: "token",
	}))

	//For checking authorization level
	in.Use(authMiddleware)

	for _, ep := range endpoints {
		switch ep.Access {
		case vsec.Super:
			configure(in, "r0/", ep)
		case vsec.Admin:
			configure(in, "r1/", ep)
		case vsec.Normal:
			configure(in, "r2/", ep)
		case vsec.Monitor:
			configure(in, "r3/", ep)
		case vsec.Public:
			configure(root, "", ep)
		}
	}
}

//Serve - start the server
func Serve(port int) (err error) {
	printConfig()
	address := fmt.Sprintf(":%d", port)
	err = e.Start(address)
	return err
}

//GetRootPath - get base URL of the configured application's REST Endpoints
func GetRootPath() string {
	return rootPath
}

func configure(grp *echo.Group, urlPrefix string, ep *Endpoint) {
	var route *echo.Route
	switch ep.Method {
	case echo.CONNECT:
		route = grp.CONNECT(urlPrefix+ep.URL, ep.Func)
	case echo.DELETE:
		route = grp.DELETE(urlPrefix+ep.URL, ep.Func)
	case echo.GET:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.HEAD:
		route = grp.HEAD(urlPrefix+ep.URL, ep.Func)
	case echo.OPTIONS:
		route = grp.OPTIONS(urlPrefix+ep.URL, ep.Func)
	case echo.PATCH:
		route = grp.PATCH(urlPrefix+ep.URL, ep.Func)
	case echo.POST:
		route = grp.POST(urlPrefix+ep.URL, ep.Func)
	case echo.PUT:
		route = grp.PUT(urlPrefix+ep.URL, ep.Func)
	case echo.TRACE:
		route = grp.TRACE(urlPrefix+ep.URL, ep.Func)
	}
	ep.Route = route
	if _, found := categories[ep.Category]; !found {
		categories[ep.Category] = make([]*Endpoint, 0, 20)
	}
	categories[ep.Category] = append(categories[ep.Category], ep)
}

func printConfig() {
	fmt.Println()
	fmt.Println("Endpoints: ")
	for category, eps := range categories {
		fmt.Printf("\t%10s\n", category)
		for _, ep := range eps {
			fmt.Printf("\t\t|-%10s - %10v - %-50s - %s\n",
				ep.Method,
				ep.Access,
				ep.Route.Path,
				ep.Comment)
		}
		fmt.Println()
	}
}
