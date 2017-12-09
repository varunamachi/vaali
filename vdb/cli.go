package vdb

import (
	"github.com/varunamachi/vaali/vcmn"
	"gopkg.in/urfave/cli.v1"
)

//mongoFlags - flags to get mongo connection options
var mongoFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "db-host",
		Value: "localhost",
		Usage: "Address of the host running mongodb",
	},
	cli.IntFlag{
		Name:  "db-port",
		Value: 27017,
		Usage: "Port on which Mongodb is listening",
	},
	cli.StringFlag{
		Name:  "db-user",
		Value: "",
		Usage: "Mongodb user name",
	},
	cli.StringFlag{
		Name:  "db-pass",
		Value: "",
		Usage: "Mongodb password for connection",
	},
}

func requireMongo(ctx *cli.Context) (err error) {
	return err
}

//MakeRequireMongo - makes ccommand to require information that is needed to
//connect to a mongodb instance
func MakeRequireMongo(cmd *cli.Command) *cli.Command {
	cmd.Flags = append(cmd.Flags, mongoFlags...)
	if cmd.Before == nil {
		cmd.Before = requireMongo
	} else {
		otherBefore := cmd.Before
		cmd.Before = func(ctx *cli.Context) (err error) {
			err = requireMongo(ctx)
			if err == nil {
				err = otherBefore(ctx)
			}
			return err
		}
	}
	return cmd
}

//GetMongoOpts - gets mongo db connection options from commandline
func GetMongoOpts(ctx *cli.Context) (opts *MongoConnOpts, ag *vcmn.ArgGetter) {
	ag = vcmn.NewArgGetter(ctx)
	host := ag.GetRequiredString("db-host")
	port := ag.GetRequiredInt("db-port")
	user := ag.GetString("db-user")
	pswd := ag.GetString("password")
	opts = &MongoConnOpts{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pswd,
	}
	return opts, ag
}

//GetCommands - get list of commands for mongdb
func GetCommands() (cmds []cli.Command) {
	cmds = []cli.Command{
		*MakeRequireMongo(
			&cli.Command{
				Name:   "test-mongo",
				Action: testMongoCmd,
			},
		),
	}
	return cmds
}

//testMongoCmd - command for testing mongodb commands
func testMongoCmd(ctx *cli.Context) (err error) {
	opts, ag := GetMongoOpts(ctx)
	if ag.Err == nil {
		err = ConnectSingle(opts)
		CloseMongoConn()
	}
	return err
}
