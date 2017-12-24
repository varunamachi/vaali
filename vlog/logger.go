package vlog

import (
	"fmt"
	"os"
	"runtime"
)

//ToString - maps level to a string
func ToString(level Level) string {
	switch level {
	case TraceLevel:
		return "[TRACE]"
	case DebugLevel:
		return "[DEBUG]"
	case InfoLevel:
		return "[ INFO]"
	case WarnLevel:
		return "[ WARN]"
	case ErrorLevel:
		return "[ERROR]"
	case FatalLevel:
		return "[FATAL]"
	}
	return "[     ]"
}

func (level Level) String() string {
	return ToString(level)
}

//InitWithOptions - initializes the logger with non default options. If you
//want default behavior, no need to call any init functions
func InitWithOptions(lc LoggerConfig) {
	lconf = lc
	if lc.LogConsole {
		lconf.Logger.RegisterWriter(NewConsoleWriter())
	}
}

//SetLevel - sets the filter level
func SetLevel(level Level) {
	lconf.FilterLevel = level
}

//GetLevel - gets the filter level
func GetLevel() (level Level) {
	return lconf.FilterLevel
}

//Trace - trace logs
func Trace(module, fmtStr string, args ...interface{}) {
	if TraceLevel >= lconf.FilterLevel {
		lconf.Logger.Log(TraceLevel, module, fmtStr, args...)
	}
}

//Debug - debug logs
func Debug(module, fmtStr string, args ...interface{}) {
	if DebugLevel >= lconf.FilterLevel {
		lconf.Logger.Log(DebugLevel, module, fmtStr, args...)
	}
}

//Info - information logs
func Info(module, fmtStr string, args ...interface{}) {
	if InfoLevel >= lconf.FilterLevel {
		lconf.Logger.Log(InfoLevel, module, fmtStr, args...)
	}
}

//Warn - warning logs
func Warn(module, fmtStr string, args ...interface{}) {
	if WarnLevel >= lconf.FilterLevel {
		lconf.Logger.Log(WarnLevel, module, fmtStr, args...)
	}
}

//Error - error logs
func Error(module, fmtStr string, args ...interface{}) {
	if ErrorLevel >= lconf.FilterLevel {
		lconf.Logger.Log(ErrorLevel, module, fmtStr, args...)
		Print(module, fmtStr, args...)
	}
}

//Fatal - error logs
func Fatal(module, fmtStr string, args ...interface{}) {
	lconf.Logger.Log(FatalLevel, module, fmtStr, args...)
	Print(module, fmtStr, args...)
	os.Exit(-1)
}

//LogError - error log
func LogError(module string, err error) error {
	if err != nil && ErrorLevel >= lconf.FilterLevel {
		_, file, line, _ := runtime.Caller(1)
		lconf.Logger.Log(ErrorLevel, module, "%v -- %s @ %d", err, file, line)
	}
	return err
}

//LogFatal - logs before exit
func LogFatal(module string, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		lconf.Logger.Log(FatalLevel, module, "%v -- %s @ %d", err, file, line)
		// Print(module, "%v", err)
		os.Exit(-1)
	}
}

//Print - prints the message on console
func Print(module, fmtStr string, args ...interface{}) {
	lconf.Logger.Log(PrintLevel, module, fmtStr, args)
	fmt.Printf(fmtStr+"\n", args...)
}
