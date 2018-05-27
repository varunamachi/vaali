package vapp

import (
	"github.com/varunamachi/vaali/vmgo"
	"github.com/varunamachi/vaali/vnet"
	"gopkg.in/urfave/cli.v1"
)

//App - the application itself
type App struct {
	cli.App
	Modules       []*Module    `json:"modules"`
	NetOptions    vnet.Options `json:"netOptions"`
	IsService     bool         `json:"isService"`
	RequiresMongo bool         `json:"requiresMongo"`
}

//ModuleConfigFunc Signature used by functions that are used to configure a
//module. Some config callbacks include - initialize, setup, reset etc
type ModuleConfigFunc func(app *App) (err error)

//Factory - wraps data type name and a function to create an instance of it
type Factory struct {
	DataType string `json:"dataType"`
	Func     vmgo.FactoryFunc
}

//Module - represents an application module
type Module struct {
	Name        string           `json:"name"`
	Description string           `json:"desc"`
	Endpoints   []*vnet.Endpoint `json:"endpoints"`
	Factories   []Factory        `json:"factory"`
	Commands    []cli.Command
	Initialize  ModuleConfigFunc
	Setup       ModuleConfigFunc
	Reset       ModuleConfigFunc
}
