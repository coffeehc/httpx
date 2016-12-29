package logger

import "flag"

var flagLoggerLevel = flag.String("logger_level", DefaultLevel, "默认日志级别")

// InitLogger 初始化日志
func InitLogger() {
	if !flag.Parsed() {
		flag.Parse()
	}
	SetDefaultLevel("/", getLevel(*flagLoggerLevel))
	Info("设置默认日志级别为:%s", *flagLoggerLevel)
}

//SetDefaultLevel 设置对应路径下默认的日志级别,可动态调整日志级别
func SetDefaultLevel(path string, level Level) {
	if path == "" {
		path = "/"
	}
	for _, filter := range filters {
		if filter.path == path {
			filter.level = level
			return
		}
	}
}
