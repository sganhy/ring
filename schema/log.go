package schema

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const levelError int8 = 1
const levelInfo int8 = 2
const levelWarning int8 = 3
const levelFatal int8 = 4
const levelDebug int8 = 5
const tagLoggerFrame = "schema.(*Schema).Log"

type log struct {
	id          int32     // {param 1}
	entryTime   time.Time // computed
	level       int8      // computed
	schemaId    int32     // {param 2}
	threadId    int32     // computed
	callSite    string    // computed
	jobId       int64     // {param 3}
	method      string    // computed
	lineNumber  int32     // computed
	message     string    // {param 4}
	description string    // {param 5}
	active      bool      // none
}

var logSchema *Schema
var logTable *Table

func (logger *log) Init(schemaId int32, disableDbLogs bool) {
	logger.schemaId = schemaId
	logger.active = !disableDbLogs
	if schemaId == 0 {
		logger.info(0, getCurrentJobId(), "Baseline Logger Initialized", "")
	} else {
		logger.info(0, getCurrentJobId(), "Logger Initialized", "")
	}
}

func initLogger(schema *Schema, table *Table) {
	logSchema = schema
	logTable = table
}

//******************************
// getters / setters
//******************************
func (logger *log) setMethod(method string) {
	if len(method) > 80 {
		logger.method = method[0:79]
	} else {
		logger.method = method
	}
}
func (logger *log) setMessage(message string) {
	if len(message) > 255 {
		logger.message = message[0:254]
	} else {
		logger.message = message
	}
}
func (logger *log) setEntryTime() {
	// truncate preserving the monotonic clock
	currTime := time.Now().UTC()
	nano := currTime.Nanosecond()
	logger.entryTime = currTime.Add(time.Nanosecond * time.Duration(-1*nano))
}
func (logger *log) setCallSite(callSite string) {

}

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************

//go:noinline
func (logger *log) warn(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelWarning, jobId, messages...)
}

//go:noinline
func (logger *log) error(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelError, jobId, messages...)
}

//go:noinline
func (logger *log) fatal(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelFatal, jobId, messages...)
}

//go:noinline
func (logger *log) debug(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelDebug, jobId, messages...)
}

//go:noinline
func (logger *log) info(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelInfo, jobId, messages...)
}

//go:noinline
func (logger *log) writePartialLog(id int32, level int8, jobId int64, messages ...interface{}) {
	var message = logger.getMessage(messages)
	var description = ""

	var newLog = logger.getNewLog(id, level, logger.schemaId, jobId, message, description)
	// am I ready to log?
	if logger.active == true && logSchema != nil && logTable != nil && logSchema.poolInitialized == true {
		logger.writeToDb(newLog)
	} else if logger.active == true {
		// connection pool not yet ready
		go logger.writeDeferToDb(newLog)
	} else {
		// print to console (sysctl)

	}
}

func (logger *log) writeToDb(newLog *log) {
	query := new(metaQuery)
	// target log table
	query.schema = logSchema
	query.table = logTable
	query.insertLog(newLog)
}

func (logger *log) writeDeferToDb(newLog *log) {
	// max try 10 times
	for i := 0; i < 10; i++ {
		fmt.Println("time.Sleep(1 * time.Second)")
		time.Sleep(1 * time.Second)
		if logger.active == true && logSchema != nil && logTable != nil && logSchema.poolInitialized == true {
			logger.writeToDb(newLog)
			break
		}
	}
}

//go:noinline
func (logger *log) getNewLog(id int32, level int8, schemaId int32, jobId int64, message string, description string) *log {
	callsite, method, lineNumber := logger.getCallerInfo()

	newLog := new(log)
	newLog.id = id
	newLog.setEntryTime()
	newLog.level = level
	newLog.threadId = 0
	newLog.callSite = callsite
	newLog.jobId = jobId
	newLog.setMethod(method)
	newLog.lineNumber = lineNumber
	newLog.setMessage(message)
	newLog.description = description

	return newLog
}

func (logger *log) getCallerInfo() (string, string, int32) {
	frame := logger.getFrame(4) // or 5

	// reach another frame?
	if strings.Contains(frame.Function, tagLoggerFrame) == true {
		frame = logger.getFrame(5)
	}
	method := frame.Function
	callSite := frame.Function
	lastIndex := strings.LastIndex(frame.Function, ".")
	if lastIndex > 0 && lastIndex+1 < len(frame.Function) {
		method = frame.Function[lastIndex+1:]
		callSite = frame.Function[:lastIndex]
	}
	return callSite, method, int32(frame.Line)
}

// get method name
func (logger *log) getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}
	return frame
}

func (logger *log) getMessage(params []interface{}) string {
	fmt.Println("(logger *log) getMessage")
	if len(params) > 0 {
		fmt.Println("(logger *log) getMessage 1 ")
		param := params[0]
		fmt.Println("(logger *log) getMessage 2 ")
		fmt.Printf("%T", param)
		switch param.(type) {
		case string:
			fmt.Println("(logger *log) getMessage 3 ")
			return param.(string)
		}

		/*
			for _, arg := range params {
				switch arg.(type) {
				case string:
					fmt.Println(arg.(string))
					return arg.(string)
				}
			}
		*/

	}
	return ""
}
