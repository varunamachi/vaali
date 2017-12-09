package vapp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/varunamachi/vaali/vnet"
	cli "gopkg.in/urfave/cli.v1"
)

//Version - represents version of the application
type Version struct {
	Major int `json:"major" bson:"major"`
	Minor int `json:"minor" bson:"minor"`
	Patch int `json:"patch" bson:"patch"`
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
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
	for _, module := range app.Modules {
		app.Commands = append(app.Commands, module.Cmds...)
		vnet.AddEndpoints(module.Endpoints...)
	}
	return app.Run(args)
}

//NewApplication - creates a new application
func NewApplication(
	name string,
	version Version,
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

// //Run - runs the application
// func (orek *OrekApp) Run(args []string) (err error) {
// 	if runtime.GOOS != "windows" {

// 	}
// 	app := cli.NewApp()
// 	app.ErrWriter = ioutil.Discard
// 	app.Name = "Orek"
// 	app.Version = "0.0.1"
// 	app.Authors = []cli.Author{
// 		cli.Author{
// 			Name:  "Varun Amachi",
// 			Email: "varunamachi@github.com",
// 		},
// 	}
// 	app.Flags = []cli.Flag{
// 		cli.StringFlag{
// 			Name:  "ds",
// 			Value: "sqlite",
// 			Usage: "Datasource name, sqlite|postgres",
// 		},
// 		cli.StringFlag{
// 			Name:  "db-path",
// 			Value: fromOrekDir("orek.db"),
// 			Usage: "Path to SQLite database [Only applicable for SQLite]",
// 		},
// 		cli.StringFlag{
// 			Name:  "ds-host",
// 			Value: "localhost",
// 			Usage: "DataBase host name [Not applicable for SqliteDataSource]",
// 		},
// 		cli.IntFlag{
// 			Name:  "ds-port",
// 			Value: 5432,
// 			Usage: "DataBase port [Not applicable for SqliteDataSource]",
// 		},
// 		cli.StringFlag{
// 			Name:  "db-name",
// 			Value: "orek",
// 			Usage: "DataBase name [Not applicable for SqliteDataSource]",
// 		},
// 		cli.StringFlag{
// 			Name:  "db-user",
// 			Value: "",
// 			Usage: "DataBase username [Not applicable for SqliteDataSource]",
// 		},
// 		cli.StringFlag{
// 			Name:  "db-password",
// 			Value: "",
// 			Usage: "Option db password for testing " +
// 				"[Not applicable for SqliteDataSource]",
// 		},
// 	}
// 	app.Before = func(ctx *cli.Context) (err error) {
// 		argetr := ArgGetter{Ctx: ctx}
// 		ds := argetr.GetRequiredString("ds")
// 		var store data.OrekDataStore
// 		if ds == "sqlite" {
// 			path := argetr.GetRequiredString("db-path")
// 			dirPath := filepath.Dir(path)
// 			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
// 				err = os.Mkdir(dirPath, 0755)
// 				olog.PrintError("Orek", err)
// 			}
// 			store, err = sqlite.Init(&sqlite.Options{
// 				Path: path,
// 			})
// 			if err == nil {
// 				data.SetStore(store)
// 				// err = data.GetStore().Init()
// 				if err != nil {
// 					olog.Fatal("Orek",
// 						"Data Store initialization failed: %v", err)
// 				} else {
// 					olog.Info("Orek", "%s Data Store initialized", store.Name())
// 				}
// 			}
// 		} else if ds == "postgres" {
// 			host := argetr.GetRequiredString("db-host")
// 			port := argetr.GetRequiredInt("db-port")
// 			dbName := argetr.GetRequiredString("db-name")
// 			user := argetr.GetRequiredString("db-user")
// 			pswd := argetr.GetString("db-password")
// 			if len(pswd) == 0 {
// 				fmt.Printf("Password for %s: ", user)
// 				var pbyte []byte
// 				pbyte, err = terminal.ReadPassword(int(syscall.Stdin))
// 				if err != nil {
// 					olog.Fatal("Orek", "Could not retrieve DB password: %v", err)
// 				} else {
// 					pswd = string(pbyte)
// 				}
// 			}
// 			olog.Print("Orek", `Postgres isnt supported yet. Here are the args
// 				Host: %s,
// 				Port: %d,
// 				DbName: %s,
// 				User: %s`, host, port, dbName, user)
// 		} else {
// 			olog.Fatal("Orek", "Unknown datasource %s requested", ds)
// 		}
// 		return err
// 	}
// 	app.Commands = make([]cli.Command, 0, 30)
// 	for _, cmdp := range orek.CommandProviders {
// 		app.Commands = append(app.Commands, cmdp.GetCommand())
// 	}
// 	err = app.Run(args)
// 	return err
// }

//OrekApp - contains command providers and runs the app
// type OrekApp struct {
// 	CommandProviders []CliCommandProvider
// }
