package vapp

import (
	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vdb"
	"github.com/varunamachi/vaali/vnet"
	"gopkg.in/urfave/cli.v1"
)

//GetCommands - gives commands related to HTTP networking
func GetCommands() []cli.Command {
	return []cli.Command{
		*vdb.MakeRequireMongo(serviceStart()),
		*vdb.MakeRequireMongo(makeAdmin()),
	}
}

func serviceStart() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Starts the HTTP service",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: 8000,
				Usage: "Port at which the service needs to serve",
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			ag := vcmn.NewArgGetter(ctx)
			port := ag.GetRequiredInt("port")
			if err = ag.Err; err == nil {
				err = vnet.Serve(port)
			}
			return err
		},
		Subcommands: []cli.Command{},
	}
}

//makeAdmin - makes an user admin
func makeAdmin() *cli.Command {
	return &cli.Command{
		Name:  "make-admin",
		Flags: []cli.Flag{},
	}
}
