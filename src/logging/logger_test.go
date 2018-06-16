package logging

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"strings"
	"testing"
)

var count = 0

func TestConsoleLogger(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	logger := &ConsoleLogger{}
	logger.SetDefaultInvokingNumber()
	expectedInvokingNumber := uint(1)
	currentInvokingNumber := logger.getInvokingNumber()
	if currentInvokingNumber != expectedInvokingNumber {
		t.Errorf("The current invoking number %d should be %d!", currentInvokingNumber, expectedInvokingNumber)
	}
	testLogger(t, logger)
}

func TestLogManager(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	logger := &LogManager{loggers: []Logger{&ConsoleLogger{invokingNumber: 2}}}
	testLogger(t, logger)
}

func testLogger(t *testing.T, logger Logger) {
	var format string
	var content string
	var logContent string

	format = ""
	logContent = "<Error>"
	content = logger.Error(logContent)
	checkContent(t, getErrorLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Errorf"
	content = logger.Errorf(format, logContent)
	checkContent(t, getErrorLogTag(), content, format, logContent)

	format = ""
	logContent = "<Errorln>"
	content = logger.Errorln(logContent)
	checkContent(t, getErrorLogTag(), content, format, logContent)

	format = ""
	logContent = "<Fatal>"
	content = logger.Fatal(logContent)
	checkContent(t, getFatalLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Fatalf"
	content = logger.Fatalf(format, logContent)
	checkContent(t, getFatalLogTag(), content, format, logContent)

	format = ""
	logContent = "<Fatalln>"
	content = logger.Fatalln(logContent)
	checkContent(t, getFatalLogTag(), content, format, logContent)

	format = ""
	logContent = "<Info>"
	content = logger.Info(logContent)
	checkContent(t, getInfoLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Infof"
	content = logger.Infof(format, logContent)
	checkContent(t, getInfoLogTag(), content, format, logContent)

	format = ""
	logContent = "<Infoln>"
	content = logger.Infoln(logContent)
	checkContent(t, getInfoLogTag(), content, format, logContent)

	format = ""
	logContent = "<Panic>"
	content = logger.Panic(logContent)
	checkContent(t, getPanicLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Panicf"
	content = logger.Panicf(format, logContent)
	checkContent(t, getPanicLogTag(), content, format, logContent)

	format = ""
	logContent = "<Panicln>"
	content = logger.Panicln(logContent)
	checkContent(t, getPanicLogTag(), content, format, logContent)

	format = ""
	logContent = "<Warn>"
	content = logger.Warn(logContent)
	checkContent(t, getWarnLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Warnf"
	content = logger.Warnf(format, logContent)
	checkContent(t, getWarnLogTag(), content, format, logContent)

	format = ""
	logContent = "<Warnln>"
	content = logger.Warnln(logContent)
	checkContent(t, getWarnLogTag(), content, format, logContent)
}

func checkContent(t *testing.T, logTag LogTag, content string, format string, logContents ...interface{}) {
	var prefixBuffer bytes.Buffer
	prefixBuffer.WriteString(logTag.Prefix())
	prefixBuffer.WriteString(" go_lib/logging.testLogger : (logger_test.go:")
	prefix := prefixBuffer.String()
	var suffixBuffer bytes.Buffer
	suffixBuffer.WriteString(") - ")
	if len(format) == 0 {
		suffixBuffer.WriteString(fmt.Sprint(logContents...))
	} else {
		suffixBuffer.WriteString(fmt.Sprintf(format, logContents...))
	}
	suffix := suffixBuffer.String()
	if !strings.HasPrefix(content, prefix) {
		t.Errorf("The content '%s' should has prefix '%s'! ", content, prefix)
	}
	if !strings.HasSuffix(content, suffix) {
		t.Errorf("The content '%s' should has suffix '%s'! ", content, suffix)
	}
}
