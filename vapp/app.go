package vapp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vnet"
	"github.com/varunamachi/vaali/vuman"
	cli "gopkg.in/urfave/cli.v1"
)

//App - the application itself
type App struct {
	cli.App
	Modules    []*Module    `json:"modules"`
	NetOptions vnet.Options `json:"netOptions"`
}

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
	fmt.Printf("Starting %s v.%v\n", app.Name, app.Version)
	vnet.AddEndpoints(vnet.GetEndpoints()...)
	vnet.AddEndpoints(vuman.GetEndpoints()...)
	app.Commands = append(app.Commands, GetCommands()...)

	for _, module := range app.Modules {
		cmds := module.CmdProvider()
		app.Commands = append(app.Commands, cmds...)
		vnet.AddEndpoints(module.Endpoints...)
	}
	vnet.InitWithOptions(app.NetOptions)
	return app.Run(args)
}

//NewDefaultApp - creates a new application with default options
func NewDefaultApp(
	name string,
	appVersion vcmn.Version,
	apiVersion string,
	authors []cli.Author,
	desc string) (app *App) {
	vlog.InitWithOptions(vlog.LoggerConfig{
		Logger:      vlog.NewDirectLogger(),
		LogConsole:  true,
		FilterLevel: vlog.TraceLevel,
		EventLogger: MongoAuditor,
	})
	app = &App{
		App: cli.App{
			Name:      name,
			Commands:  make([]cli.Command, 0, 100),
			Version:   appVersion.String(),
			Authors:   authors,
			Usage:     desc,
			ErrWriter: ioutil.Discard,
		},
		NetOptions: vnet.Options{
			RootName:      name,
			APIVersion:    apiVersion,
			Authenticator: vuman.MongoAuthenticator,
			Authorizer:    nil,
		},
		Modules: make([]*Module, 0, 10),
	}
	return app
}

//NewAppWithOptions - creates app with non default options
func NewAppWithOptions( /*****/ ) (app *App) {
	return app
}
