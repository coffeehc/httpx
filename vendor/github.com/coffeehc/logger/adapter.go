package logger

import (
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"errors"
)

var errorLoggerClose = errors.New("logger is close")

type internalWaiter struct {
	_filter *logFilter
}

func (waiter *internalWaiter) Write(p []byte) (n int, err error) {
	if isStop {
		return -1, errorLoggerClose
	}
	if len(p) == 0 {
		return 0, nil
	}
	_, file, line, ok := runtime.Caller(4)
	var lineInfo = "-:0"
	if ok {
		index := strings.Index(file, "/src/") + 4
		lineInfo = file[index:] + ":" + strconv.Itoa(line)
	}
	log := &logContent{time.Now(), waiter._filter.level, lineInfo, string(p)}
	waiter._filter.cache <- log
	return len(p), nil
}

//CreatLoggerAdapter 创建 *log.Logger 的适配器
func CreatLoggerAdapter(level Level, timeFormat string, format string, out io.Writer) *log.Logger {
	filter := newFilter(level, "/", timeFormat, format, out)
	go filter.run()
	return log.New(&internalWaiter{filter}, "", 0)
}
