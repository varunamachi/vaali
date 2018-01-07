package vapp

import (
	"errors"
	"fmt"
	"time"

	"github.com/varunamachi/vaali/vlog"

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
		Name:  "create-super",
		Usage: "Create a super user",
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
		Name:  "setup",
		Usage: "Sets up the application",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "super-id",
				Usage: "Super user ID",
			},
			cli.StringFlag{
				Name:  "super-pw",
				Usage: "Super user password",
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			vapp := GetVApp(ctx)
			if vapp != nil {
				ag := vcmn.NewArgGetter(ctx)
				superID := ag.GetRequiredString("super-id")
				superPW := ag.GetOptionalString("super-pw")
				if err = ag.Err; err == nil {
					defer func() {
						vlog.LogEvent(
							"setup app",
							superID,
							err != nil,
							err,
							nil)
					}()
					if len(superPW) == 0 {
						superPW = vcmn.AskPassword("Super-user Password")
					}
					var user *vsec.User
					user, err = vnet.DoLogin(superID, superPW)
					if err != nil {
						err = fmt.Errorf(
							"Failed to authenticate super user: %v",
							err)
						return err
					}
					if user.Auth != vsec.Super {
						err = errors.New(
							"User forcing reset is not a super user")
					}
					err = vapp.Setup()
				}
			} else {
				err = errors.New("V App not properly initialized")
			}
			return vlog.LogError("App", err)
		},
	}
}

func resetCmd() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "Sets up the application",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "super-id",
				Usage: "Super user ID",
			},
			cli.StringFlag{
				Name:  "super-pw",
				Usage: "Super user password",
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			vapp := GetVApp(ctx)
			if vapp != nil {
				ag := vcmn.NewArgGetter(ctx)
				superID := ag.GetRequiredString("super-id")
				superPW := ag.GetOptionalString("super-pw")
				if err = ag.Err; err == nil {
					defer func() {
						vlog.LogEvent(
							"rerset app",
							superID,
							err != nil,
							err,
							nil)
					}()
					if len(superPW) == 0 {
						superPW = vcmn.AskPassword("Super-user Password")
					}
					var user *vsec.User
					user, err = vnet.DoLogin(superID, superPW)
					if err != nil {
						err = fmt.Errorf(
							"Failed to authenticate super user: %v",
							err)
						return err
					}
					if user.Auth != vsec.Super {
						err = errors.New(
							"User forcing reset is not a super user")
					}
					err = vapp.Reset()
				}
			} else {
				err = errors.New("V App not properly initialized")
			}
			return vlog.LogError("App", err)
		},
	}
}

func overridePassword() *cli.Command {
	return &cli.Command{
		Name: "force-pw-reset",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "id",
				Usage: "Unique ID of the user",
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "New password",
			},
			cli.StringFlag{
				Name:  "super-id",
				Usage: "Super user ID",
			},
			cli.StringFlag{
				Name:  "super-pw",
				Usage: "Super user password",
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			ag := vcmn.NewArgGetter(ctx)
			id := ag.GetRequiredString("id")
			pw := ag.GetOptionalString("password")
			superID := ag.GetRequiredString("super-id")
			superPW := ag.GetOptionalString("super-pw")
			defer func() { vlog.LogError("App", err) }()
			if err = ag.Err; err == nil {
				defer func() {
					vlog.LogEvent(
						"force set password",
						superID,
						err != nil,
						err,
						vlog.M{
							"super": superID,
							"user":  id,
						})
				}()
				if len(pw) == 0 {
					pw = vcmn.AskPassword("New Password")
				}
				if len(superPW) == 0 {
					superPW = vcmn.AskPassword("Super-user Password")
				}
				var user *vsec.User
				user, err = vnet.DoLogin(superID, superPW)
				if err != nil {
					err = fmt.Errorf("Failed to authenticate super user: %v",
						err)
					return err
				}
				if user.Auth != vsec.Super {
					err = errors.New("User forcing reset is not a super user")
				}
				err = vuman.SetPassword(id, pw)
				if err == nil {
					vlog.Info("App", "Password for %s successfully reset", id)
				}
			}
			return err
		},
	}
}

//GetVApp - gets instance of vapp.App which is stored inside cli.App.Metadata
func GetVApp(ctx *cli.Context) (vapp *App) {
	metadata := ctx.App.Metadata
	vi, found := metadata["vapp"]
	if found {
		vapp, _ = vi.(*App)
	}
	return vapp
}
