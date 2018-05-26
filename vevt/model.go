package vevt

import "time"

//EventFilterModel - model for creating event filters for fields
type EventFilterModel struct {
	UserNames  []string `json:"userNames" bson:"userNames" sql:"userNames"`
	EventTypes []string `json:"eventTypes" bson:"eventTypes" sql:"eventTypes"`
}

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

//EventHandler - handles application events for audit purposes
type EventHandler interface {
	LogEvent(event *Event)
}

// func defaultAuditor(event *Event) {
// 	if event.Success {
// 		fmt.Printf("Event:Info - %s BY %s", event.Op, event.UserID)
// 	} else {
// 		fmt.Printf("Event:Error - %s BY %s", event.Op, event.UserID)
// 	}

// }

// //LogEvent - logs an event using the registered audit function
// func LogEvent(
// 	op string,
// 	userID string,
// 	userName string,
// 	success bool,
// 	err string,
// 	data interface{}) {
// 	lconf.EventLogger(&Event{
// 		Op:       op,
// 		UserID:   userID,
// 		UserName: userName,
// 		Success:  success,
// 		Error:    err,
// 		Time:     time.Now(),
// 		Data:     data,
// 	})
// }
