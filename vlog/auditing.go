package vlog

import "fmt"
import "time"

func defaultAuditor(event *Event) {
	if event.Success {
		fmt.Printf("Event:Info - %s BY %s", event.Op, event.UserID)
	} else {
		fmt.Printf("Event:Error - %s BY %s", event.Op, event.UserID)
	}

}

//LogEvent - logs an event using the registered audit function
func LogEvent(
	op string,
	userID string,
	userName string,
	success bool,
	err error,
	data interface{}) {
	estr := ""
	if err != nil {
		estr = err.Error()
	}
	lconf.EventLogger(&Event{
		Op:       op,
		UserID:   userID,
		UserName: userName,
		Success:  success,
		Error:    estr,
		Time:     time.Now(),
		Data:     data,
	})
}
