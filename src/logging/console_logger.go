package logging

import (
	"log"
)

type ConsoleLogger struct {
	position Position
}

func (logger *ConsoleLogger) GetPosition() Position {
	return logger.position
}

func (logger *ConsoleLogger) SetPosition(pos Position) {
	logger.position = pos
}

func (logger *ConsoleLogger) Error(v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Errorf(format string, v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Errorln(v ...interface{}) string {
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) Fatal(v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Fatalf(format string, v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Fatalln(v ...interface{}) string {
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) Info(v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Infof(format string, v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Infoln(v ...interface{}) string {
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) Panic(v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Panicf(format string, v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Panicln(v ...interface{}) string {
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}

func (logger *ConsoleLogger) Warn(v ...interface{}) string {
	content := generateLogContent(getWarnLogTag(), logger.GetPosition(), "", v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Warnf(format string, v ...interface{}) string {
	content := generateLogContent(getWarnLogTag(), logger.GetPosition(), format, v...)
	log.Print(content)
	return content
}

func (logger *ConsoleLogger) Warnln(v ...interface{}) string {
	content := generateLogContent(getWarnLogTag(), logger.GetPosition(), "", v...)
	log.Println(content)
	return content
}
