package vmgo

import (
	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vlog"
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
	ag := vcmn.NewArgGetter(ctx)
	dbHost := ag.GetRequiredString("db-host")
	dbPort := ag.GetRequiredInt("db-port")
	dbUser := ag.GetOptionalString("db-user")
	dbPassword := ""
	if len(dbUser) != 0 {
		dbPassword = ag.GetRequiredSecret("db-pass")
	}
	if err = ag.Err; err == nil {
		err = ConnectSingle(&MongoConnOpts{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
		})
	}
	if err != nil {
		vlog.LogFatal("DB:Mongo", err)
	}
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
// func GetMongoOpts(ctx *cli.Context) (opts *MongoConnOpts, ag *vcmn.ArgGetter) {
// 	ag = vcmn.NewArgGetter(ctx)
// 	host := ag.GetRequiredString("db-host")
// 	port := ag.GetRequiredInt("db-port")
// 	user := ag.GetString("db-user")
// 	pswd := ag.GetString("password")
// 	opts = &MongoConnOpts{
// 		Host:     host,
// 		Port:     port,
// 		User:     user,
// 		Password: pswd,
// 	}
// 	return opts, ag
// }

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
	CloseMongoConn()
	return err
}
