package logging

type LogManager struct {
	loggers []Logger
}

func (logger *LogManager) GetPosition() Position {
	return POSITION_SINGLE
}

func (logger *LogManager) SetPosition(pos Position) {}

func (self *LogManager) Error(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Error(v...)
	}
	return content
}

func (self *LogManager) Errorf(format string, v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Errorf(format, v...)
	}
	return content
}

func (self *LogManager) Errorln(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Errorln(v...)
	}
	return content
}

func (self *LogManager) Fatal(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Fatal(v...)
	}
	return content
}

func (self *LogManager) Fatalf(format string, v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Fatalf(format, v...)
	}
	return content
}

func (self *LogManager) Fatalln(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Fatalln(v...)
	}
	return content
}

func (self *LogManager) Info(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Info(v...)
	}
	return content
}

func (self *LogManager) Infof(format string, v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Infof(format, v...)
	}
	return content
}

func (self *LogManager) Infoln(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Infoln(v...)
	}
	return content
}

func (self *LogManager) Panic(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Panic(v...)
	}
	return content
}

func (self *LogManager) Panicf(format string, v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Panicf(format, v...)
	}
	return content
}

func (self *LogManager) Panicln(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Panicln(v...)
	}
	return content
}

func (self *LogManager) Warn(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Warn(v...)
	}
	return content
}

func (self *LogManager) Warnf(format string, v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Warnf(format, v...)
	}
	return content
}

func (self *LogManager) Warnln(v ...interface{}) string {
	var content string
	for _, logger := range self.loggers {
		content = logger.Warnln(v...)
	}
	return content
}
