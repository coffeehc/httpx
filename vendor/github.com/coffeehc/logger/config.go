// config

/*
 使用goconfig来获取配置，格式如下：
-
 level: debug
 package_path: /
 adapter: console
-
 level: error
 package_path: /
 adapter: file
 log_path: /logs/box/box.log
 rotate: 3
 #备份策略：size or time  or default
 rotate_policy: time
 #备份范围：如果策略是time则表示时间间隔N分钟，如果是size则表示每个日志的最大大小(MB)
 rotate_scope: 10
*/

package logger

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//日志配置
type LoggerConfig struct {
	//配置上下文
	Context string `yaml:"context"`
	//指定的日志处理机
	Appenders []LoggerAppender `yaml:"appenders"`
}

const (
	//适配定义,控制台方式
	Adapter_Console = "console"
	//适配定义,文件方式
	Adapter_File = "file"
)

//通过flag来指定日志文件路径,没有指定则查找当前目录下./conf/log.yml
var _loggerConf *string = flag.String("logger", getDefaultLog(), "日志文件路径")
var default_Config *LoggerConfig = &LoggerConfig{Context: "Default", Appenders: []LoggerAppender{{Level: LOGGER_DEFAULT_LEVEL, Package_path: "/", Adapter: "console"}}}

//获取默认的日志配置文件,路径为程序当前目录下./conf/log.yml
func getDefaultLog() string {
	file, _ := exec.LookPath(os.Args[0])
	filePath, _ := filepath.Abs(file)
	return path.Join(filepath.Dir(filePath), "conf/log.yml")
}

//加载日志配置,如果指定了-loggerConf参数,则加载这个参数指定的配置文件,如果没有则使用默认的配置
func loadLoggerConfig(loggerConf string) {
	if len(filters) > 0 {
		for _, filter := range filters {
			filter.clear()
		}
	}
	filters = make([]*logFilter, 0)
	conf := parseConfile(loggerConf)
	if conf == nil || len(conf.Appenders) == 0 {
		//fmt.Println("没有指定配置文件,服务将使用默认配置")
		conf = default_Config
	}
	for _, appender := range conf.Appenders {
		AddAppender(appender)
	}

}

//解析配置
func parseConfile(loggerConf string) *LoggerConfig {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("解析配置文件%s出错:%s\n", loggerConf, r)
		}
	}()
	//log.Printf("加载日志配置文件:%s\n", loggerConf)
	data, err := ioutil.ReadFile(loggerConf)
	if err != nil {
		//log.Printf("[警告]加载日志配置文件错误:%s\n", err)
	} else {
		conf := new(LoggerConfig)
		err = yaml.Unmarshal(data, conf)
		if err != nil {
			log.Printf("加载日志配置文件失败:%s\n", err)
		}
		return conf
	}
	return nil
}
