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

var ERR_CLOSE = errors.New("logger is close")

type internalWaiter struct {
	_filter *logFilter
}

func (this *internalWaiter) Write(p []byte) (n int, err error) {
	if isStop {
		return -1, ERR_CLOSE
	}
	if len(p) == 0 {
		return 0, nil
	}
	_, file, line, ok := runtime.Caller(4)
	var lineInfo string = "-:0"
	if ok {
		index := strings.Index(file, "/src/") + 4
		lineInfo = file[index:] + ":" + strconv.Itoa(line)
	}
	log := &logContent{time.Now(), this._filter.level, lineInfo, string(p)}
	this._filter.cache <- log
	return len(p), nil
}

func CreatLoggerAdapter(level byte, timeFormat string, format string, out io.Writer) *log.Logger {
	filter := newFilter(level, "/", timeFormat, format, out)
	go filter.run()
	return log.New(&internalWaiter{filter}, "", 0)
}
