package vnet

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/varunamachi/vaali/vsec"
)

var categories = make(map[string][]*Endpoint)
var endpoints = make([]*Endpoint, 0, 200)
var e = echo.New()
var accessPos = 0

// var groups = make

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

//Init - initializes all the registered endpoints
func Init(rootName, apiVersion string) {
	rootPath := rootName + "/api/v" + apiVersion + "/"
	accessPos = len(rootPath) + len("in/")
	root := e.Group(rootPath)
	in := root.Group("in/")
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

func configure(grp *echo.Group, urlPrefix string, ep *Endpoint) {
	var route *echo.Route
	switch ep.Method {
	case echo.CONNECT:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.DELETE:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.GET:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.HEAD:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.OPTIONS:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.PATCH:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.POST:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.PUT:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	case echo.TRACE:
		route = grp.GET(urlPrefix+ep.URL, ep.Func)
	}
	ep.Route = route
	eps, found := categories[ep.Category]
	if !found {
		eps = make([]*Endpoint, 0, 20)
		categories[ep.Category] = eps
	}
	eps = append(eps, ep)
}

func printConfig() {
	for category, eps := range categories {
		fmt.Println(category)
		for _, ep := range eps {
			fmt.Printf("\t%10s - %10v - %s",
				ep.Method,
				ep.Access,
				ep.Route.Path)
		}
	}
}
