package vapp

import (
	"time"

	"github.com/varunamachi/vaali/vuman"

	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vdb"
	"github.com/varunamachi/vaali/vnet"
	"github.com/varunamachi/vaali/vsec"
	"gopkg.in/urfave/cli.v1"
)

//GetCommands - gives commands related to HTTP networking
func GetCommands() []cli.Command {
	return []cli.Command{
		*vdb.MakeRequireMongo(serviceStartCmd()),
		*vdb.MakeRequireMongo(createUserCmd()),
		*vdb.MakeRequireMongo(setupCmd()),
		*vdb.MakeRequireMongo(resetCmd()),
	}
}

func serviceStartCmd() *cli.Command {
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
// func makeAdmin() *cli.Command {
// 	return &cli.Command{
// 		Name:  "make-admin",
// 		Flags: []cli.Flag{},
// 	}
// }

func createUserCmd() *cli.Command {
	return &cli.Command{
		Name: "create-super",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "id",
				Usage: "Unique ID of the user",
			},
			cli.StringFlag{
				Name:  "email",
				Usage: "Email of the password",
			},
			cli.StringFlag{
				Name:  "first",
				Usage: "First name of the user",
			},
			cli.StringFlag{
				Name:  "last",
				Usage: "Last name of the user",
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			ag := vcmn.NewArgGetter(ctx)
			id := ag.GetRequiredString("id")
			email := ag.GetRequiredString("email")
			first := ag.GetString("first")
			last := ag.GetString("last")
			if err = ag.Err; err == nil {
				one := vcmn.AskPassword("Password")
				two := vcmn.AskPassword("Confirm")
				if one == two {
					user := vsec.User{
						ID:        id,
						Email:     email,
						Auth:      vsec.Super,
						FirstName: first,
						LastName:  last,
						Created:   time.Now(),
						Modified:  time.Now(),
						Props:     make(map[string]string),
						PwdExpiry: time.Now().AddDate(1, 0, 0),
					}
					err = vuman.CreateFirstSuperUser(&user, one)
				}
			}
			return err
		},
	}
}

func setupCmd() *cli.Command {
	return &cli.Command{
		Name: "setup",
	}
}

func resetCmd() *cli.Command {
	return &cli.Command{
		Name: "reset",
	}
}
