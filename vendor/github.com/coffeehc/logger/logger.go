package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//日志格式定义
const (
	FormatTime     = "%T" //显示时间
	FormatLevel    = "%L" //显示级别
	FormatCodeInfo = "%C" //显示代码详情
	FormatMessage  = "%M" //显示日志消息
)

//Logger 日志接口
type Logger interface {
	Trace(format string, v ...interface{}) string
	Debug(format string, v ...interface{}) string
	Info(format string, v ...interface{}) string
	Warn(format string, v ...interface{}) string
	Error(format string, v ...interface{}) string
}

//GetLogger 获取一个日志对象,主要用于和第三方包做适配
func GetLogger() Logger {
	return loggercopy
}

var loggercopy = &_logger{}

type _logger struct {
}

func (*_logger) Trace(format string, v ...interface{}) string {
	return output(LoggerLevelTrace, fmt.Sprintf(format, v...), 3)
}

func (*_logger) Debug(format string, v ...interface{}) string {
	return output(LoggerLevelDebug, fmt.Sprintf(format, v...), 3)
}

func (*_logger) Info(format string, v ...interface{}) string {
	return output(LoggerLevelInfo, fmt.Sprintf(format, v...), 3)
}

func (*_logger) Warn(format string, v ...interface{}) string {
	return output(LoggerLevelWarn, fmt.Sprintf(format, v...), 3)
}

func (*_logger) Error(format string, v ...interface{}) string {
	return output(LoggerLevelError, fmt.Sprintf(format, v...), 3)
}

var (
	//LoggerDefaultLevel 默认的日志级别
	LoggerDefaultLevel = "error"
	//LoggerDefaultBufsize 日志缓冲区默认大小
	LoggerDefaultBufsize = 1024
	//LoggerDefaultTimeout 默认超时时间
	LoggerDefaultTimeout = time.Second * 1
)

// Level 日志级别
type Level byte

const (
	// LoggerLevelError  Error level
	LoggerLevelError Level = 1 << 0
	// LoggerLevelWarn Warn level
	LoggerLevelWarn Level = 1<<1 | LoggerLevelError
	//LoggerLevelInfo Info level
	LoggerLevelInfo Level = 1<<2 | LoggerLevelWarn
	// LoggerLevelDebug Debug level
	LoggerLevelDebug Level = 1<<3 | LoggerLevelInfo
	// LoggerLevelTrace Trace level
	LoggerLevelTrace Level = 1<<4 | LoggerLevelDebug

	// LoggerTimeformatSecond 显示格式为2006-01-02 15:04:05
	LoggerTimeformatSecond string = "2006-01-02 15:04:05"
	// LoggerTimeformatNanosecond 显示格式为2006-01-02 15:04:05.999999999
	LoggerTimeformatNanosecond string = "2006-01-02 15:04:05.999999999"
	// LoggerTimeformatAll 显示格式为2006-01-02 15:04:05.999999999 -0700 UTC
	LoggerTimeformatAll string = "2006-01-02 15:04:05.999999999 -0700 UTC"
	loggerCodeDepth     int    = 2
)

func getLevelStr(level Level) string {
	switch level {
	case LoggerLevelError:
		return "ERROR"
	case LoggerLevelWarn:
		return "WARN"
	case LoggerLevelInfo:
		return "INFO"
	case LoggerLevelDebug:
		return "DEBUG"
	case LoggerLevelTrace:
		return "TRACE"
	default:
		return "DEBUG"
	}
}

//Flusher 日志持久化接口
type Flusher interface {
	Flush() error
}

//日志拦截器定义
type logFilter struct {
	level       Level       //拦截级别
	path        string      //拦截路径
	timeFormat  string      //时间戳格式
	format      []logFormat //日志格式
	out         io.Writer
	cache       chan *logContent
	stop        chan bool
	filterClose chan bool
}

//判断是否需要过滤器处理
func (filter *logFilter) canSave(level Level, lineInfo string) bool {
	return filter.level&level == level && strings.HasPrefix(lineInfo, filter.path)
}

func (filter *logFilter) save(buf *bytes.Buffer, content *logContent) {
	buf.Reset()
	for _, f := range filter.format {
		buf.Write(f.data)
		switch f.appendContent {
		case FormatTime:
			buf.WriteString(content.appendTime.Format(filter.timeFormat))
		case FormatLevel:
			buf.WriteString(getLevelStr(content.level))
		case FormatCodeInfo:
			buf.WriteString(content.lineInfo[1:])
		case FormatMessage:
			buf.WriteString(content.content)
		default:
		}
	}
	if content.content[len(content.content)-1] != '\n' {
		buf.WriteByte('\n')
	}
	buf.WriteTo(filter.out)
}

//过滤器后台输出goruntine
func (filter *logFilter) run() {
	timeOut := time.Millisecond * 500
	timer := time.NewTimer(timeOut)
	buf := bytes.NewBuffer(nil)
	stop := false
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("logger error :%s", err)
		}
	}()
	for {
		select {
		case content := <-filter.cache:
			if filter.canSave(content.level, content.lineInfo) {
				filter.save(buf, content)
			}
		case <-timer.C:
			if v, ok := filter.out.(Flusher); ok {
				v.Flush()
			}
			if stop {
				goto CLOSE
			}
		case <-filter.stop:
			stop = true
			if len(filter.cache) == 0 {
				goto CLOSE
			}
			timeOut = time.Millisecond * 100
		}
		timer.Reset(timeOut)
	}
CLOSE:
	close(filter.filterClose)
}

func (filter *logFilter) clear() {
	close(filter.stop)
	<-filter.filterClose
}

var (
	filters []*logFilter
	isStop  = false
)

//启动日志
func init() {
	loadLoggerConfig(*_loggerConf)
}

func addStdOutFilter(level Level, path string, timeFormat string, format string) {
	AddFileter(level, path, timeFormat, format, os.Stdout)
}

//ClearFilter 清空过滤器,主要用于自定义处理日志
func ClearFilter() {
	filters = make([]*logFilter, 0)
}

//AddFileter 添加日志过滤器,参数说明:级别,包路径,时间格式,Writer接口
func AddFileter(level Level, path string, timeFormat string, format string, out io.Writer) {
	filter := newFilter(level, path, timeFormat, format, out)
	filters = append(filters, filter)
	go filter.run()
}

func newFilter(level Level, path string, timeFormat string, format string, out io.Writer) *logFilter {
	if timeFormat == "" {
		timeFormat = LoggerTimeformatSecond
	}
	if path == "" {
		path = "/"
	}
	if timeFormat == "" {
		timeFormat = LoggerTimeformatSecond
	}
	if format == "" {
		format = "%T %L %C %M"
	}
	if out == nil {
		panic("拦截器输出不能为空")
	}
	filter := new(logFilter)
	filter.level = level
	filter.path = path
	filter.timeFormat = timeFormat
	filter.out = out
	filter.format = parseFormat(format)
	filter.cache = make(chan *logContent, 200)
	filter.stop = make(chan bool, 10)
	filter.filterClose = make(chan bool)
	return filter
}

type logFormat struct {
	data          []byte
	appendContent string
}

func parseFormat(format string) []logFormat {
	fs := make([]logFormat, 0)
	fc := []byte(format)
	start := 0
	i := 0
	for ; i < len(fc); i++ {
		b := fc[i]
		if b == '%' {
			switch fc[i+1] {
			case 'T':
				fs = append(fs, logFormat{data: fc[start:i], appendContent: FormatTime})
				i++
				start = i + 1
				continue
			case 'L':
				fs = append(fs, logFormat{data: fc[start:i], appendContent: FormatLevel})
				i++
				start = i + 1
				break
			case 'C':
				fs = append(fs, logFormat{data: fc[start:i], appendContent: FormatCodeInfo})
				i++
				start = i + 1
				break
			case 'M':
				fs = append(fs, logFormat{data: fc[start:i], appendContent: FormatMessage})
				i++
				start = i + 1
				break
			}
		}
	}
	if i > start {
		fs = append(fs, logFormat{data: fc[start : i-1]})
	}
	return fs
}

type logContent struct {
	appendTime time.Time
	level      Level
	lineInfo   string
	content    string
}

//WaitToClose called then the application end
func WaitToClose() {
	isStop = true
	for _, filter := range filters {
		if filter != nil {
			filter.clear()
		}
	}
}

// real out implement
func output(logLevel Level, content string, codeLevel int) string {
	if isStop {
		return ""
	}
	if len(content) == 0 {
		return ""
	}
	_, file, line, ok := runtime.Caller(codeLevel)
	lineInfo := "-:0"
	if ok {
		index := strings.Index(file, "/src/") + 4
		lineInfo = file[index:] + ":" + strconv.Itoa(line)
	}
	log := &logContent{time.Now(), logLevel, lineInfo, content}
	for _, filter := range filters {
		if filter != nil {
			filter.cache <- log
		}
	}
	return content
}

//Printf  logger Printer
func Printf(logLevel Level, codeLevel int, format string, v ...interface{}) string {
	return output(logLevel, fmt.Sprintf(format, v...), codeLevel)
}

//Trace print trace log
func Trace(format string, v ...interface{}) string {
	return output(LoggerLevelTrace, fmt.Sprintf(format, v...), loggerCodeDepth)
}

//Debug print debug log
func Debug(format string, v ...interface{}) string {
	return output(LoggerLevelDebug, fmt.Sprintf(format, v...), loggerCodeDepth)
}

//Info print Info log
func Info(format string, v ...interface{}) string {
	return output(LoggerLevelInfo, fmt.Sprintf(format, v...), loggerCodeDepth)
}

//Warn print Warn log
func Warn(format string, v ...interface{}) string {
	return output(LoggerLevelWarn, fmt.Sprintf(format, v...), loggerCodeDepth)
}

// Error print Error log
func Error(format string, v ...interface{}) string {
	return output(LoggerLevelError, fmt.Sprintf(format, v...), loggerCodeDepth)
}
