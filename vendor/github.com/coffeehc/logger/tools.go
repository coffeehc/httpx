// tool
package logger

import "flag"

var flag_logger_level = flag.String("logger_level", LOGGER_DEFAULT_LEVEL, "默认日志级别")

//初始化日志
func InitLogger() {
	if !flag.Parsed() {
		flag.Parse()
	}
	SetDefaultLevel("/", getLevel(*flag_logger_level))
	Info("设置默认日志级别为:%s", *flag_logger_level)
}

//设置对应路径下默认的日志级别,可动态调整日志级别
func SetDefaultLevel(path string, level byte) {
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
