package vevt

import (
	"fmt"
	"time"

	"github.com/varunamachi/vaali/vmgo"
)

// //EventFilterModel - model for creating event filters for fields
// type EventFilterModel struct {
// 	UserNames  []string `json:"userNames" bson:"userNames" sql:"userNames"`
// 	EventTypes []string `json:"eventTypes" bson:"eventTypes" sql:"eventTypes"`
// }

// //NewEventFilterModel - creates a new event filter model
// func NewEventFilterModel() EventFilterModel {
// 	return EventFilterModel{
// 		UserNames:  make([]string, 0, 1000),
// 		EventTypes: make([]string, 0, 100),
// 	}
// }

//Event - represents a event initiated by a user while performing an operation
type Event struct {
	Op       string      `json:"op" bson:"op"`
	UserID   string      `json:"userID" bson:"userID"`
	UserName string      `json:"userName" bson:"userName"`
	Success  bool        `json:"success" bson:"success"`
	Error    string      `json:"error" bson:"error"`
	Time     time.Time   `json:"time" bson:"time"`
	Data     interface{} `json:"data" bson:"data"`
}

//EventAuditor - handles application events for audit purposes
type EventAuditor interface {
	//LogEvent - logs given event into storage
	LogEvent(event *Event)

	//GetEvents - retrieves event entries based on filters
	GetEvents(
		offset, limit int, filter *vmgo.Filter) (
		total int, events []*Event, err error)

	//CreateIndices - creates mongoDB indeces for tables used for event logs
	CreateIndices() (err error)

	//CleanData - cleans event related data from database
	CleanData() (err error)
}

//NoOpAuditor - doesnt do anything, it's a dummy auditor
type NoOpAuditor struct{}

//LogEvent - logs event to console
func (n *NoOpAuditor) LogEvent(event *Event) {
	if event.Success {
		fmt.Printf("Event:Info - %s BY %s", event.Op, event.UserID)
	} else {
		fmt.Printf("Event:Error - %s BY %s", event.Op, event.UserID)
	}
}

//GetEvents - gives an empty list of events
func (n *NoOpAuditor) GetEvents(
	offset, limit int, filter *vmgo.Filter) (
	total int, events []*Event, err error) {
	return total, events, err
}

//CreateIndices - creates nothing
func (n *NoOpAuditor) CreateIndices() (err error) { return err }

//CleanData - there's nothing to clean
func (n *NoOpAuditor) CleanData() (err error) { return err }
