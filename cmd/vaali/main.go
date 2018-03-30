package main

import (
	"os"

	"github.com/varunamachi/vaali/vapp"
	"github.com/varunamachi/vaali/vcmn"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := vapp.NewWebApp(
		"vaali",
		vcmn.Version{
			Major: 0,
			Minor: 0,
			Patch: 0,
		},
		"0",
		[]cli.Author{
			cli.Author{
				Name: "Varuna Amachi",
			},
		},
		"Default app for Vaali",
	)
	app.Exec(os.Args)
}
