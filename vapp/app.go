package vapp

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/varunamachi/vaali/vevt"
	"github.com/varunamachi/vaali/vsec"

	"github.com/varunamachi/vaali/vmgo"

	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vlog"
	"github.com/varunamachi/vaali/vnet"
	"github.com/varunamachi/vaali/vuman"
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
	if app.IsService {
		vnet.AddEndpoints(vnet.GetEndpoints()...)
		vnet.AddEndpoints(vuman.GetEndpoints()...)
		vnet.AddEndpoints(getEndpoints()...)
	}
	if app.RequiresMongo {
		app.Commands = append(app.Commands, GetCommands(app)...)
	}

	for _, module := range app.Modules {
		if module.Initialize != nil {
			err = module.Initialize(app)
			if err != nil {
				vlog.Error("App", "Failed to initialize module %s",
					module.Name)
				break
			}
		}
		if module.Commands != nil {
			app.Commands = append(app.Commands, module.Commands...)
		}
		if module.Factories != nil {
			for _, fc := range module.Factories {
				vmgo.RegisterFactory(fc.DataType, fc.Func)
			}
		}
		if app.IsService {
			vnet.AddEndpoints(module.Endpoints...)
		}
	}
	if err == nil {
		vnet.InitWithOptions(app.NetOptions)
		err = app.Run(args)
	}
	return err
}

//NewWebApp - creates a new web application with default options
func NewWebApp(
	name string,
	appVersion vcmn.Version,
	apiVersion string,
	authors []cli.Author,
	requiresMongo bool, desc string) (app *App) {
	var store vsec.UserStorage
	var auditor vevt.EventAuditor
	store = &vuman.MongoStorage{}
	auditor = &vevt.MongoAuditor{}
	authr := vuman.MongoAuthenticator
	if !requiresMongo {
		store = &vuman.PGStorage{}
		auditor = &vevt.PGAuditor{}
	}
	vevt.SetEventAuditor(auditor)
	vlog.InitWithOptions(vlog.LoggerConfig{
		Logger:      vlog.NewDirectLogger(),
		LogConsole:  true,
		FilterLevel: vlog.TraceLevel,
	})

	vcmn.LoadConfig(name)
	app = &App{
		UserStorage:   store,
		IsService:     true,
		RequiresMongo: true,
		App: cli.App{
			Name:      name,
			Commands:  make([]cli.Command, 0, 100),
			Version:   appVersion.String(),
			Authors:   authors,
			Usage:     desc,
			ErrWriter: ioutil.Discard,
			Metadata:  map[string]interface{}{},
		},
		NetOptions: vnet.Options{
			RootName:      name,
			APIVersion:    apiVersion,
			Authenticator: authr,
			Authorizer:    nil,
		},
		Modules: make([]*Module, 0, 10),
	}
	app.Metadata["vapp"] = app
	return app
}

//NewSimpleApp - an app that is not a service and does not use mongodb
func NewSimpleApp(
	name string,
	appVersion vcmn.Version,
	apiVersion string,
	authors []cli.Author,
	desc string) (app *App) {
	vlog.InitWithOptions(vlog.LoggerConfig{
		Logger:      vlog.NewDirectLogger(),
		LogConsole:  true,
		FilterLevel: vlog.TraceLevel,
	})
	vevt.SetEventAuditor(&vevt.NoOpAuditor{})
	vcmn.LoadConfig(name)
	app = &App{
		IsService:     false,
		RequiresMongo: false,
		App: cli.App{
			Name:      name,
			Commands:  make([]cli.Command, 0, 100),
			Version:   appVersion.String(),
			Authors:   authors,
			Usage:     desc,
			ErrWriter: ioutil.Discard,
			Metadata:  map[string]interface{}{},
		},
		Modules: make([]*Module, 0, 10),
	}
	app.Metadata["vapp"] = app
	return app
}

//Setup - sets up the application and the registered module. This is not
//initialization and needs to be called when app/module configuration changes.
//This is the place where mongoDB indices are expected to be created.
func (app *App) Setup() (err error) {
	if app.RequiresMongo {

		err = vuman.CreateIndices()
		if err != nil {
			vlog.Error("App",
				"Failed to create Mongo indeces for user storage")
			return err
		}
		vlog.Info("App", "Created indeces for user storage")
		err = vevt.GetAuditor().CreateIndices()
		if err != nil {
			vlog.Error("App",
				"Failed to create Mongo indeces for event storage")
			return err
		}
		vlog.Info("App", "Created indeces for user storage")
	}
	for _, module := range app.Modules {
		if module.Setup != nil {
			err = module.Setup(app)
			if err != nil {
				vlog.Error("App", "Failed to set module %s up",
					module.Name)
				return err
			}
			vlog.Info("App", "Configured module %s", module.Name)
		}
	}
	if err != nil {
		err = errors.New("Failed to set the application up")
	} else {
		vlog.Info("App", "Application setup complete")
	}
	return err
}

//Reset - resets the application and module configuration and data.
//USE WITH CAUTION
func (app *App) Reset() (err error) {
	if app.RequiresMongo {
		err = vuman.CleanData()
		if err != nil {
			vlog.Error("App", "Failed to reset user storage")
		}
	}
	for _, module := range app.Modules {
		if module.Setup != nil {
			err = module.Reset(app)
			if err != nil {
				vlog.Error("App", "Failed to reset module %s",
					module.Name)
			} else {
				vlog.Info("App", "Reset module %s", module.Name)
			}
		}
	}
	if err != nil {
		err = errors.New("Failed to reset application")
	}
	return err
}

//NewAppWithOptions - creates app with non default options
func NewAppWithOptions( /*****/ ) (app *App) {
	return app
}
