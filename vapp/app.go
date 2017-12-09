package vapp

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vnet"
	cli "gopkg.in/urfave/cli.v1"
)

//FromAppDir - gives a absolute path from a path relative to
//app directory
func (app *App) FromAppDir(relPath string) (abs string) {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("APPDATA")
	}
	return filepath.Join(home, "."+app.Name, relPath)
}

//AddModule - registers a module with the app
func (app *App) AddModule(module *Module) {
	app.Modules = append(app.Modules, module)
}

//Exec - runs the applications
func (app *App) Exec(args []string) (err error) {
	for _, module := range app.Modules {
		cmds := module.CmdProvider()
		app.Commands = append(app.Commands, cmds...)
		vnet.AddEndpoints(module.Endpoints...)
	}
	return app.Run(args)
}

//NewApplication - creates a new application
func NewApplication(
	name string,
	version vcmn.Version,
	authors []cli.Author,
	desc string,
) (app *App) {
	app = &App{
		App: cli.App{
			Name:      name,
			Commands:  make([]cli.Command, 0, 100),
			Version:   version.String(),
			Authors:   authors,
			Usage:     desc,
			ErrWriter: ioutil.Discard,
		},
		Modules: make([]*Module, 0, 10),
	}
	return app
}
