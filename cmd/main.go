package main

import (
	"github.com/varunamachi/vaali/vapp"
	"github.com/varunamachi/vaali/vlog"
)

func main() {
	vlog.InitWithOptions(vlog.LoggerConfig{
		Logger:      vlog.NewDirectLogger(),
		LogConsole:  true,
		FilterLevel: vlog.InfoLevel,
		EventLogger: vapp.MongoAuditor,
	})
}
