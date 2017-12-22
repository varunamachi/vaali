package vnet

import (
	"github.com/varunamachi/vaali/vcmn"
	"gopkg.in/urfave/cli.v1"
)

//GetCommands - gives commands related to HTTP networking
func GetCommands() []cli.Command {
	return []cli.Command{
		cli.Command{
			Name: "service",
			Subcommands: []cli.Command{
				serviceStart(),
			},
		},
	}
}

func serviceStart() cli.Command {
	return cli.Command{
		Name: "start",
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
				Serve(port)
			}
			return err
		},
	}
}
