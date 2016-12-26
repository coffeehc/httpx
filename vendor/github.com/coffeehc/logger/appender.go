package logger

import (
	"fmt"
	"strings"
	"time"
)

// Appender 日志处理机
type Appender struct {
	Level        string `yaml:"level"`         //日志级别
	PackagePath  string `yaml:"package_path"`  //日志路径
	Adapter      string `yaml:"adapter"`       //适配器,console,file两种
	Rotate       int    `yaml:"rotate"`        //日志切割个数
	RotatePolicy string `yaml:"rotate_policy"` //切割策略,time or size or  default
	RotateScope  int64  `yaml:"rotate_scope"`  //切割范围:如果按时间切割则表示的n分钟,如果是size这表示的是文件大小MB
	LogPath      string `yaml:"log_path"`      //如果适配器使用的file则用来指定文件路径
	Timeformat   string `yaml:"timeformat"`    //指定时间输出格式.默认:"2006-01-02 15:04:05"
	Format       string `yaml:"format"`        //日志格式
}

// AddAppender 添加日志配置,支持Console与File两种方式
func AddAppender(appender *Appender) {
	switch appender.Adapter {
	case AdapterFile:
		addFileAppender(appender)
		break
	case AdapterConsole:
		addConsoleAppender(appender)
		break
	default:
		fmt.Printf("不能识别的日志适配器:%s", appender.Adapter)
	}
}

//添加console的日志配置
func addConsoleAppender(appender *Appender) {
	addStdOutFilter(getLevel(appender.Level), appender.PackagePath, appender.Timeformat, appender.Format)
}

//添加文件系统的日志配置
func addFileAppender(appender *Appender) {
	rotatePolicy := strings.ToLower(appender.RotatePolicy)
	switch rotatePolicy {
	case "time":
		addFileFilterForTime(getLevel(appender.Level), appender.PackagePath, appender.LogPath, time.Minute*time.Duration(appender.RotateScope), appender.Rotate, appender.Timeformat, appender.Format)
		return
	case "size":
		addFileFilterForSize(getLevel(appender.Level), appender.PackagePath, appender.LogPath, appender.RotateScope*1048576, appender.Rotate, appender.Timeformat, appender.Format)
		return
	default:
		addFileFilterForDefualt(getLevel(appender.Level), appender.PackagePath, appender.LogPath, appender.Timeformat, appender.Format)
		return
	}
}

func getLevel(level string) Level {
	level = strings.ToLower(level)
	switch level {
	case "trace":
		return LoggerLevelTrace
	case "debug":
		return LoggerLevelDebug
	case "info":
		return LoggerLevelInfo
	case "warn":
		return LoggerLevelWarn
	case "error":
		return LoggerLevelError
	default:
		return LoggerLevelDebug
	}
}
