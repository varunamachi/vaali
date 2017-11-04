package vlog

import (
	"fmt"
)

//DirectLogger - logger that writes directly to all registered writers
type DirectLogger struct {
	writers map[string]Writer
}

//NewDirectLogger - creates a new DirectLogger instace
func NewDirectLogger() *DirectLogger {
	return &DirectLogger{
		writers: make(map[string]Writer),
	}
}

//Log - logs a message with given level and module
func (dl *DirectLogger) Log(level Level,
	module string,
	fmtstr string,
	args ...interface{}) {
	if level == PrintLevel {
		return
	}
	fmtstr = ToString(level) + " [" + module + "] " + fmtstr
	msg := fmt.Sprintf(fmtstr, args...)
	for _, writer := range dl.writers {
		if writer.IsEnabled() {
			writer.Write(msg)
		}
	}
}

//RegisterWriter - registers a writer
func (dl *DirectLogger) RegisterWriter(writer Writer) {
	if writer != nil {
		dl.writers[writer.UniqueID()] = writer
	}
}

//RemoveWriter - removes a writer with given ID
func (dl *DirectLogger) RemoveWriter(uniqueID string) {
	delete(dl.writers, uniqueID)
}

//GetWriter - gives the writer with given ID
func (dl *DirectLogger) GetWriter(uniqueID string) (writer Writer) {
	return dl.writers[uniqueID]
}
