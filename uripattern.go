package httpx

import (
	"fmt"
	"regexp"
	"strings"

	"git.xiagaogao.com/coffee/boot/logs"
	"go.uber.org/zap"
)

//URIPatternMatcher a URI matcher interface
type URIPatternMatcher interface {
	match(uri string) bool
}

type servletStyleURIPatternMatcher struct {
	pattern     string
	patternKind int //0:PREFIX 1:SUFFIX 2:LITERAL
}

func newServletStyleURIPatternMatcher(uriPattern string,logger *zap.Logger) URIPatternMatcher {
	matcher := new(servletStyleURIPatternMatcher)
	if strings.HasPrefix(uriPattern, "*") {
		matcher.pattern = string([]byte(uriPattern)[1:])
		matcher.patternKind = 0
	} else if strings.HasSuffix(uriPattern, "*") {
		matcher.pattern = string([]byte(uriPattern)[:len(uriPattern)-1])
		matcher.patternKind = 1
	} else {
		matcher.pattern = uriPattern
		matcher.patternKind = 2
	}
	return URIPatternMatcher(matcher)
}
func (matcher *servletStyleURIPatternMatcher) match(uri string) bool {
	if uri == "" {
		return false
	}
	switch matcher.patternKind {
	case 0:
		return strings.HasSuffix(uri, matcher.pattern)
	case 1:
		return strings.HasPrefix(uri, matcher.pattern)
	default:
		return matcher.pattern == uri
	}
}

type regexURIPatternMatcher struct {
	pattern *regexp.Regexp
}

func newRegexURIPatternMatcher(uriPattern string,logger *zap.Logger) URIPatternMatcher {
	matcher := new(regexURIPatternMatcher)
	var err error
	matcher.pattern, err = regexp.Compile(uriPattern)
	if err != nil {
		logger.Error(fmt.Sprintf("编译正则表达式[%s]异常", uriPattern), logs.F_Error(err))
		return nil
	}
	return URIPatternMatcher(matcher)
}

func (matcher *regexURIPatternMatcher) match(uri string) bool {
	return uri != "" && matcher.pattern.MatchString(uri)
}
