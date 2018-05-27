package vevt

import "time"

var eventAuditor EventAuditor

//SetEventAuditor - sets the event auditor
func SetEventAuditor(auditor EventAuditor) {
	eventAuditor = auditor
}

//SetEventAuditor - gets the event auditor
func GetAuditor() EventAuditor {
	return eventAuditor
}

//LogEvent - logs an event using the registered audit function
func LogEvent(
	op string,
	userID string,
	userName string,
	success bool,
	err string,
	data interface{}) {
	eventAuditor.LogEvent(&Event{
		Op:       op,
		UserID:   userID,
		UserName: userName,
		Success:  success,
		Error:    err,
		Time:     time.Now(),
		Data:     data,
	})
}
