package logging

const (
	ERROR_LOG_KEY = "ERROR"
	FATAL_LOG_KEY = "FATAL"
	INFO_LOG_KEY  = "INFO"
	PANIC_LOG_KEY = "PANIC"
	WARN_LOG_KEY  = "WARN"
)

type LogTag struct {
	name   string
	prefix string
}

func (self *LogTag) Name() string {
	return self.name
}

func (self *LogTag) Prefix() string {
	return self.prefix
}

var logTagMap map[string]LogTag = map[string]LogTag{
	ERROR_LOG_KEY: LogTag{name: ERROR_LOG_KEY, prefix: "[" + ERROR_LOG_KEY + "]"},
	FATAL_LOG_KEY: LogTag{name: FATAL_LOG_KEY, prefix: "[" + FATAL_LOG_KEY + "]"},
	INFO_LOG_KEY:  LogTag{name: INFO_LOG_KEY, prefix: "[" + INFO_LOG_KEY + "]"},
	PANIC_LOG_KEY: LogTag{name: PANIC_LOG_KEY, prefix: "[" + PANIC_LOG_KEY + "]"},
	WARN_LOG_KEY:  LogTag{name: WARN_LOG_KEY, prefix: "[" + WARN_LOG_KEY + "]"},
}

func getErrorLogTag() LogTag {
	return logTagMap[ERROR_LOG_KEY]
}

func getFatalLogTag() LogTag {
	return logTagMap[FATAL_LOG_KEY]
}

func getInfoLogTag() LogTag {
	return logTagMap[INFO_LOG_KEY]
}

func getPanicLogTag() LogTag {
	return logTagMap[PANIC_LOG_KEY]
}

func getWarnLogTag() LogTag {
	return logTagMap[WARN_LOG_KEY]
}
