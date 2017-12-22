package vapp

import (
	"github.com/varunamachi/vaali/vnet"
	"gopkg.in/urfave/cli.v1"
)

//CmdProvider - gives all the commands for a module
type CmdProvider func() []cli.Command

//Module - represents an application module
type Module struct {
	Name        string           `json:"name"`
	Description string           `json:"desc"`
	Endpoints   []*vnet.Endpoint `json:"endpoints"`
	CmdProvider CmdProvider      `json:"cmdProvider"`
}
