package vapp

import (
	"errors"
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
	vnet.AddEndpoints(getEndpoints()...)
	app.Commands = append(app.Commands, GetCommands()...)

	for _, module := range app.Modules {
		if module.Initialize != nil {
			err = module.Initialize(app)
			if err != nil {
				vlog.Error("App", "Failed to initialize module %s",
					module.Name)
				break
			}
		}
		cmds := module.CmdProvider()
		app.Commands = append(app.Commands, cmds...)
		vnet.AddEndpoints(module.Endpoints...)
	}
	if err == nil {
		vnet.InitWithOptions(app.NetOptions)
		err = app.Run(args)
	}
	return err
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
	vcmn.LoadConfig(name)
	// pstr := vcmn.GetConfigDef("smtpPort", "586")
	// port, e := strconv.Atoi(pstr)
	// if e != nil {
	// 	port = 586
	// }
	// ecfg := vnet.EmailConfig{
	// 	From:     vcmn.GetConfig("appEMail"),
	// 	Password: vcmn.GetConfig("appEMailPassword"),
	// 	SMTPHost: vcmn.GetConfig("smtpHost"),
	// 	SMTPPort: port,
	// }
	// printConfig()
	app = &App{
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
			Authenticator: vuman.MongoAuthenticator,
			Authorizer:    nil,
			EmailConfig:   nil,
			// EmailConfig:   ecfg,

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
	err = vuman.CreateIndices()
	if err != nil {
		vlog.Error("App",
			"Failed to create Mongo indeces for U-Man collections")
		return err
	}
	vlog.Info("App", "Created indeces for User Management collections")
	err = CreateIndices()
	if err != nil {
		vlog.Error("App",
			"Failed to create Mongo indeces for applications collections")
		return err
	}
	vlog.Info("App", "Created indeces for Application collections")
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
	err = vuman.CleanData()
	if err != nil {
		vlog.Error("App", "Failed to reset U-Man data")
	}
	for _, module := range app.Modules {
		if module.Setup != nil {
			err = module.Setup(app)
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
