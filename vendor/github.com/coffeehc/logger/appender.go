// appender
package logger

import (
	"fmt"
	"strings"
	"time"
)

//日志处理机
type LoggerAppender struct {
	Level         string `yaml:"level"`         //日志级别
	Package_path  string `yaml:"package_path"`  //日志路径
	Adapter       string `yaml:"adapter"`       //适配器,console,file两种
	Rotate        int    `yaml:"rotate"`        //日志切割个数
	Rotate_policy string `yaml:"rotate_policy"` //切割策略,time or size or  default
	Rotate_scope  int64  `yaml:"rotate_scope"`  //切割范围:如果按时间切割则表示的n分钟,如果是size这表示的是文件大小MB
	Log_path      string `yaml:"log_path"`      //如果适配器使用的file则用来指定文件路径
	Timeformat    string `yaml:"timeformat"`    //指定时间输出格式.默认:"2006-01-02 15:04:05"
	Format        string `yaml:"format"`        //日志格式
}

//添加日志配置,支持Console与File两种方式
func AddAppender(appender LoggerAppender) {
	switch appender.Adapter {
	case Adapter_File:
		addFileAppender(appender)
		break
	case Adapter_Console:
		addConsoleAppender(appender)
		break
	default:
		fmt.Printf("不能识别的日志适配器:%s", appender.Adapter)
	}
}

//添加console的日志配置
func addConsoleAppender(appender LoggerAppender) {
	addStdOutFilter(getLevel(appender.Level), appender.Package_path, appender.Timeformat, appender.Format)
}

//添加文件系统的日志配置
func addFileAppender(appender LoggerAppender) {
	rotatePolicy := strings.ToLower(appender.Rotate_policy)
	switch rotatePolicy {
	case "time":
		addFileFilterForTime(getLevel(appender.Level), appender.Package_path, appender.Log_path, time.Minute*time.Duration(appender.Rotate_scope), appender.Rotate, appender.Timeformat, appender.Format)
		return
	case "size":
		addFileFilterForSize(getLevel(appender.Level), appender.Package_path, appender.Log_path, appender.Rotate_scope*1048576, appender.Rotate, appender.Timeformat, appender.Format)
		return
	default:
		addFileFilterForDefualt(getLevel(appender.Level), appender.Package_path, appender.Log_path, appender.Timeformat, appender.Format)
		return
	}
}

func getLevel(level string) byte {
	level = strings.ToLower(level)
	switch level {
	case "trace":
		return LOGGER_LEVEL_TRACE
	case "debug":
		return LOGGER_LEVEL_DEBUG
	case "info":
		return LOGGER_LEVEL_INFO
	case "warn":
		return LOGGER_LEVEL_WARN
	case "error":
		return LOGGER_LEVEL_ERROR
	default:
		return LOGGER_LEVEL_DEBUG
	}
}
