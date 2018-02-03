package vapp

import (
	"github.com/varunamachi/vaali/vnet"
	"gopkg.in/urfave/cli.v1"
)

//CmdProvider - gives all the commands for a module
type CmdProvider func() []cli.Command

//ModuleConfigFunc Signature used by functions that are used to configure a
//module. Some config callbacks include - initialize, setup, reset etc
type ModuleConfigFunc func(app *App) (err error)

//Module - represents an application module
type Module struct {
	Name        string           `json:"name"`
	Description string           `json:"desc"`
	Endpoints   []*vnet.Endpoint `json:"endpoints"`
	CmdProvider CmdProvider
	Initialize  ModuleConfigFunc
	Setup       ModuleConfigFunc
	Reset       ModuleConfigFunc
}

//EventFilterModel - model for creating event filters for fields
type EventFilterModel struct {
	UserNames  []string `json:"userNames" bson:"userNames"`
	EventTypes []string `json:"eventTypes" bson:"eventTypes"`
}

//NewEventFilterModel - creates a new event filter model
func NewEventFilterModel() EventFilterModel {
	return EventFilterModel{
		UserNames:  make([]string, 0, 1000),
		EventTypes: make([]string, 0, 100),
	}
}
