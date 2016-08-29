package web

import "github.com/coffeehc/logger"

type httpLogger struct {
}

func (this httpLogger) Printf(format string, args ...interface{}) {
	logger.Printf(logger.LOGGER_LEVEL_DEBUG, logger.LOGGER_CODE_DEPTH+1, format, args)
}
