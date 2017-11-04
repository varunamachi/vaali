package log

import "fmt"
import "sync"

//AsyncLogger - logger that uses goroutine for dispatching
type AsyncLogger struct {
	sync.Mutex
	writers map[string]Writer
}

//NewAsyncLogger - creates a new DirectLogger instace
func NewAsyncLogger() *AsyncLogger {
	return &AsyncLogger{
		writers: make(map[string]Writer),
	}
}

//Log - logs a message with given level and module
func (al *AsyncLogger) Log(level Level,
	module string,
	fmtstr string,
	args ...interface{}) {
	if level == PrintLevel {
		return
	}
	go func() {
		fmtstr = ToString(level) + " [" + module + "] " + fmtstr
		msg := fmt.Sprintf(fmtstr, args...)
		al.Lock()
		for _, writer := range al.writers {
			if writer.IsEnabled() {
				writer.Write(msg)
			}
		}
		al.Unlock()
	}()
}

//RegisterWriter - registers a writer
func (al *AsyncLogger) RegisterWriter(writer Writer) {
	if writer != nil {
		al.Lock()
		al.writers[writer.UniqueID()] = writer
		al.Unlock()
	}
}

//RemoveWriter - removes a writer with given ID
func (al *AsyncLogger) RemoveWriter(uniqueID string) {
	al.Lock()
	delete(al.writers, uniqueID)
	al.Unlock()
}

//GetWriter - gives the writer with given ID
func (al *AsyncLogger) GetWriter(uniqueID string) (writer Writer) {
	al.Lock()
	l := al.writers[uniqueID]
	al.Unlock()
	return l
}
