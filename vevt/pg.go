package vevt

import (
	"github.com/varunamachi/vaali/vcmn"
	"github.com/varunamachi/vaali/vlog"
)

//PGAuditor - handles application events and stores them in postgres for audit
// purposes
type PGAuditor struct{}

//LogEvent - stores event logs into database
func (m *PGAuditor) LogEvent(event *Event) {

}

//GetEvents - retrieves event entries based on filters
func (m *PGAuditor) GetEvents(offset, limit int, filter *vcmn.Filter) (
	total int, events []*Event, err error) {
	return total, events, vlog.LogError("App:Event:PG", err)
}

//CreateIndices - creates mongoDB indeces for tables used for event logs
func (m *PGAuditor) CreateIndices() (err error) {
	return vlog.LogError("App:Event:PG", err)
}

//CleanData - cleans event related data from database
func (m *PGAuditor) CleanData() (err error) {
	return vlog.LogError("App:Event:PG", err)
}
