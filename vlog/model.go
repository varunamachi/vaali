package vlog

//Level - gives log level
type Level int

//M - map of string to interface for representing event data
type M map[string]interface{}

const (
	//TraceLevel - low level debug message
	TraceLevel Level = 1

	//DebugLevel - a debug message
	DebugLevel Level = 2

	//InfoLevel - information message
	InfoLevel Level = 3

	//WarnLevel - warning message
	WarnLevel Level = 4

	//ErrorLevel - error message
	ErrorLevel Level = 5

	//FatalLevel - fatal messages
	FatalLevel Level = 6

	//PrintLevel - prints a output message
	PrintLevel Level = 7
)

//Writer - interface that takes a message and writes it based on
//the implementation
type Writer interface {
	UniqueID() string
	Write(message string)
	Enable(value bool)
	IsEnabled() (value bool)
}

//Logger - interface that defines a logger implementation
type Logger interface {
	//Log - logs a message with given level and module
	Log(level Level,
		module string,
		fmtstr string,
		args ...interface{})

	//RegisterWriter - registers a writer
	RegisterWriter(writer Writer)

	//RemoveWriter - removes a writer with given ID
	RemoveWriter(uniqueID string)

	//GetWriter - gives the writer with given ID
	GetWriter(uniqueID string) (writer Writer)
}

//LoggerConfig - configuration that is used to initialize the logger
type LoggerConfig struct {
	Logger      Logger
	LogConsole  bool
	FilterLevel Level
}

var lconf = LoggerConfig{
	Logger:      NewDirectLogger(),
	LogConsole:  false,
	FilterLevel: InfoLevel,
}
