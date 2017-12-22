package vuman

import (
	"gopkg.in/urfave/cli.v1"
)

//GetCommands - gives commands related to user management
func GetCommands() []cli.Command {
	return []cli.Command{
		makeAdmin(),
	}
}

//makeAdmin - makes an user admin
func makeAdmin() cli.Command {
	return cli.Command{
		Name:  "make-admin",
		Flags: []cli.Flag{},
	}
}
