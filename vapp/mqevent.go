package vapp

import (
	"github.com/varunamachi/vaali/vdb"
	"github.com/varunamachi/vaali/vlog"
)

//MongoAuditor - stores event logs into database
func MongoAuditor(event *vlog.Event) {
	conn := vdb.DefaultMongoConn()
	defer conn.Close()
	conn.C("events").Insert(event)
}
