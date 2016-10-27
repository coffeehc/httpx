logger
[![GoDoc](https://godoc.org/github.com/coffeehc/logger?status.png)](http://godoc.org/github.com/coffeehc/logger)
==

###获取方式
```
go get github.com/coffeehc/logger
```

1.0.1 更新日志
增加了 log.Logger的适配器方法
使用方式 
```go
    CreatLoggerAdapter(level byte, timeFormat string, format string, out io.Writer) *log.Logger
```
用于在某些可以注入 log.Logger的地方使用,输出仍然能适配为本项目格式,Writer 同样可以使用内置的FileWrite

=======================

###使用方式
coffeehc/logger 是一个基础的日志框架,提供扩展的开放logFilter接口,使用者可以自己定义何种级别的logger发布到对应的io.Writer中

日志级别定义了5个:

1.	trace
2.	debug
3.	info
4.	warn
5.	error

#### API说明

```go
    func Trace(format string, v ...interface{}) string
    func Debug(format string, v ...interface{}) string
    func Warn(format string, v ...interface{}) string
    func Info(format string, v ...interface{}) string
    func Error(format string, v ...interface{}) string
```

#### 编码方式定义Appender,用于程序自主定义日志

```go
    func AddAppender(appender LoggerAppender)

    type LoggerAppender struct {
        Level         string `yaml:"level"`         //日志级别
        Package_path  string `yaml:"package_path"`  //日志路径
        Adapter       string `yaml:"adapter"`       //适配器,console,file两种
        Rotate        int    `yaml:"rotate"`        //日志切割个数
        Rotate_policy string `yaml:"rotate_policy"` //切割策略,time or size or  default
        Rotate_scope  int64  `yaml:"rotate_scope"`  //切割范围:如果按时间切割则表示的n分钟,如果是size这表示的是文件大小MB
        Log_path      string `yaml:"log_path"`      //如果适配器使用的file则用来指定文件路径
        Timeformat    string `yaml:"timeformat"`    //日志格式
        Format        string `yaml:"format"`
}
```

#### 自定义更低级别的Filter

```go

    func AddFileter(level byte, path string, timeFormat string, format string, out io.Writer)
```
只需要将out实现为任意想输出的方式,tcp,http,db等都可以

#### 配置说明
使用配置的方式(yaml语法),配置文件内容如下:

```json
context: Default
appenders:
 -
  level: debug
  package_path: /
  adapter: console
  #使用golang自己的timeFormat
  timeformat: 2006-01-02 15:04:05
  format: %T %L %C %M
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
```
系统默认会取-loggerConf参数的值来加载配置文件,如果没有指定则使用debug对所有的包路径下的日志打印到控制台

2015-4-29
1. AddAppender用于自己使用编程方式来定义日志,其实也可以用底层的Filter接口来扩展会更灵活

2015-5-15
在配置中加入了format的参数设置,提供四种标记来组合日志,标记说明如下:

> 1. %T:时间标记,会与timeformat配合使用
> 2. %L:日志级别,这会输出相应的日志级别
> 3. %C:代码信息,这包括包文件描述和日志在第几行打印
> 4. %M:这个就是需要打印的具体日志内容

支持在程序运行目录下查找conf/log.yml文件作为默认的日志配置,如果不指定-loggerConf的话,可以直接将配置文件放在这下面,程序启动的时候可以直接读取配置

###TODO
1. 暂不支持TCP方式存储日志,以后看情况再提供,只要实现io.Writer的接口就可以了,自己动手,丰衣足食
