package schema

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

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

const (
	levelError       int8  = 1
	levelInfo        int8  = 2
	levelWarning     int8  = 3
	levelFatal       int8  = 4
	levelDebug       int8  = 5
	schemaNotDefined int32 = -2
)

var (
	logSchema      *Schema
	logTable       *Table
	metaTableExist bool
)

func init() {
	metaTableExist = false
}

func (logger *log) Init(schemaId int32, jobId int64, disableDbLogs bool) {
	logger.schemaId = schemaId
	logger.active = !disableDbLogs

	if disableDbLogs == true {
		// return here to avoid crash during unitests
		return
	}

	if schemaId != schemaNotDefined {
		// get metas chema to fetch current JobId
		if schemaId == 0 {
			logger.Info(0, 0, "Baseline logger initialized", "")
		} else if jobId == 0 {
			logger.Info(0, 0, "Logger initialized", "")
		}
	}
}

// get @log table + its schema
func initLogger(schema *Schema, table *Table) {
	logSchema = schema
	logTable = table
}

//******************************
// getters / setters
//******************************
func (logger *log) getId() int32 {
	return logger.id
}
func (logger *log) getEntryTime() time.Time {
	return logger.entryTime
}
func (logger *log) getLevel() int8 {
	return logger.level
}
func (logger *log) getSchemaId() int32 {
	return logger.schemaId
}
func (logger *log) getThreadId() int32 {
	return logger.threadId
}
func (logger *log) getClassSite() string {
	return logger.callSite
}
func (logger *log) getJobId() int64 {
	return logger.jobId
}
func (logger *log) getMethod() string {
	return logger.method
}
func (logger *log) getLineNumber() int32 {
	return logger.lineNumber
}
func (logger *log) getMessage() string {
	return logger.message
}
func (logger *log) getDescription() string {
	return logger.description
}

func (logger *log) setMethod(method string) {
	if len(method) > 80 {
		logger.method = method[0:80]
	} else {
		logger.method = method
	}
}
func (logger *log) setMessage(message string) {
	if len(message) > 255 {
		logger.message = message[0:255]
	} else {
		logger.message = message
	}
}
func (logger *log) setEntryTime() {
	// truncate preserving the monotonic clock
	currTime := time.Now().UTC()
	nano := currTime.Nanosecond() % 1000000
	// truncate nano seconds
	logger.entryTime = currTime.Add(time.Duration(-nano) * time.Nanosecond)
}
func (logger *log) setCallSite(callSite string) {
	if len(callSite) > 255 {
		logger.callSite = callSite[0:255]
	} else {
		logger.callSite = callSite
	}
}

// is meta table present in the database
func (logger *log) isMetaTable(exist bool) {
	metaTableExist = exist
}

func (logger *log) setSchemaId(id int32) {
	logger.schemaId = id
}

//******************************
// public methods
//******************************
//go:noinline
func (logger *log) Warn(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelWarning, jobId, messages...)
}

//go:noinline
func (logger *log) Error(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelError, jobId, messages...)
}

//go:noinline
func (logger *log) Fatal(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelFatal, jobId, messages...)
}

//go:noinline
func (logger *log) Debug(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelDebug, jobId, messages...)
}

//go:noinline
func (logger *log) Info(id int32, jobId int64, messages ...interface{}) {
	logger.writePartialLog(id, levelInfo, jobId, messages...)
}

//******************************
// private methods
//******************************
//go:noinline
func (logger *log) writePartialLog(id int32, level int8, jobId int64, messages ...interface{}) {
	var message = logger.extractMessage(messages)
	var description = logger.extractDescription(messages)

	var newLog = logger.getNewLog(id, level, logger.schemaId, jobId, message, description)
	// am I ready to log?
	if logger.active == true && logSchema != nil && logTable != nil && logSchema.poolInitialized == true && metaTableExist == true {
		//fmt.Println("logger.writeToDb(newLog) ==> " + description)
		logger.writeToDb(newLog)
	} else if logger.active == true {
		// connection pool not yet ready
		//fmt.Println("go logger.writeDeferToDb(newLog) ==> " + description)
		go logger.writeDeferToDb(newLog)
	} else {
		//do nothing?
		//fmt.Println(" not log ==> " + description)
		//fmt.Println(message)
	}
}

func (logger *log) writeToDb(newLog *log) {
	query := new(metaQuery)
	query.Init(logSchema, logTable)
	query.insertLog(newLog)
}

func (logger *log) writeDeferToDb(newLog *log) {
	// max try 1 minute

	for i := 0; i < 30; i++ {
		//fmt.Println("writeDeferToDb")
		time.Sleep(2 * time.Second)
		if logSchema != nil && logTable != nil && logSchema.poolInitialized == true && metaTableExist == true {
			// retrieve current jobId
			var metaSchema = GetSchemaByName(metaSchemaName)
			newLog.jobId = metaSchema.getJobIdNextValue()
			newLog.jobId = initialJobId
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
	newLog.threadId = logger.getCurrentThreadId()
	newLog.callSite = callsite
	newLog.jobId = jobId
	newLog.schemaId = schemaId
	newLog.setMethod(method)
	newLog.lineNumber = lineNumber
	newLog.setMessage(message)
	newLog.description = description

	return newLog
}

func (logger *log) getCallerInfo() (string, string, int32) {
	frame := logger.getFrame(4) // or 5

	// reach another frame?
	/*tagLoggerFrame   string = "schema.(*Schema).Log"
	if strings.Contains(frame.Function, tagLoggerFrame) == true {
		frame = logger.getFrame(5)
	}
	*/
	method := frame.Function
	callSite := frame.Function
	lastIndex := strings.LastIndex(frame.Function, ".")
	if lastIndex > 0 && lastIndex+1 < len(frame.Function) {
		method = frame.Function[lastIndex+1:]
		callSite = frame.Function[:lastIndex]
	}
	return callSite, method, int32(frame.Line)
}

func (logger *log) getCurrentThreadId() int32 {
	return int32(os.Getpid() & 2147483647)
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

func (logger *log) extractMessage(params []interface{}) string {
	if len(params) > 0 {
		param := params[0]
		switch param.(type) {
		case string:
			return param.(string)
		case error:
			var err = param.(error)
			return err.Error()
		}
	}
	return ""
}

/*TEST ==>
func LogTest(id int32, jobId int64, messages ...interface{}) {
	logTest := new(log)
	logTest.Init(schemaNotDefined, false)
	logTest.writePartialLog(id, levelError, jobId, messages)
}
*/

func (logger *log) extractDescription(params []interface{}) string {
	if len(params) > 1 {
		param := params[1]
		switch param.(type) {
		case string:
			return param.(string)
		}
	}
	if len(params) > 0 {
		param := params[0]
		switch param.(type) {
		case error:
			return logger.getStackTrace()
		}
	}
	return ""
}

func (logger *log) getStackTrace() string {
	stackBuf := make([]uintptr, 50)
	length := runtime.Callers(3, stackBuf[:])
	stack := stackBuf[:length]
	frameCount := 0

	trace := ""
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if frameCount > 1 {
			trace = trace + fmt.Sprintf("File: %s, Line: %d. Function: %s \n", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
		frameCount++
	}
	return trace
}
