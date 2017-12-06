package vapp

import (
	"gopkg.in/urfave/cli.v1"
)

//Module - represents an application module
type Module struct {
	Name        string      `json:"name"`
	Description string      `json:"desc"`
	Endpoints   []string    `json:"endpoints"`
	Cmd         cli.Command `json:"command"`
}

//App - the application itself
type App struct {
	cli.App
	Modules []*Module `json:"modules"`
}
