package main

import (
	"fmt"
	"os"

	"github.com/varunamachi/vaali/vlog"

	"github.com/varunamachi/vaali/vnet"

	"github.com/varunamachi/vaali/vapp"
	"github.com/varunamachi/vaali/vcmn"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := vapp.NewSimpleApp(
		"vlclient",
		vcmn.Version{
			Major: 0,
			Minor: 0,
			Patch: 1,
		},
		"0",
		[]cli.Author{
			cli.Author{
				Name: "Varuna Amachi",
			},
		},
		"Simple vaali client",
	)
	app.Commands = append(app.Commands, cli.Command{
		Name: "login",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "userID",
				Value: "",
				Usage: "Vaali user ID",
			},
			cli.StringFlag{
				Name:   "password",
				Value:  "",
				Usage:  "Vaali password",
				Hidden: true,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			ag := vcmn.NewArgGetter(ctx)
			userID := ag.GetRequiredString("userID")
			password := ag.GetOptionalString("password")
			if len(password) == 0 {
				password = vcmn.AskPassword("Password")
			}
			if err = ag.Err; err == nil {
				c := vnet.NewClient("http://localhost:8000", "vaali", "v0")
				err = c.Login(userID, password)
				if err == nil {
					fmt.Println("Login successful. User: ")
					vcmn.DumpJSON(c.User)
				}
			}
			return vlog.LogError("Client", err)
		},
	})
	app.Exec(os.Args)
}
