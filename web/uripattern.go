// pathmatch
package web

import (
	"regexp"
	"strings"

	"github.com/coffeehc/logger"
)

type UriPatternMatcher interface {
	match(uri string) bool
}

type servletStyleUriPatternMatcher struct {
	pattern     string
	patternKind int //0:PREFIX 1:SUFFIX 2:LITERAL
}

func newServletStyleUriPatternMatcher(uriPattern string) UriPatternMatcher {
	matcher := new(servletStyleUriPatternMatcher)
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
	return UriPatternMatcher(matcher)
}
func (this *servletStyleUriPatternMatcher) match(uri string) bool {
	if uri == "" {
		return false
	}
	switch this.patternKind {
	case 0:
		return strings.HasSuffix(uri, this.pattern)
	case 1:
		return strings.HasPrefix(uri, this.pattern)
	default:
		return this.pattern == uri
	}
}

type regexUriPatternMatcher struct {
	pattern *regexp.Regexp
}

func newRegexUriPatternMatcher(uriPattern string) UriPatternMatcher {
	matcher := new(regexUriPatternMatcher)
	var err error
	matcher.pattern, err = regexp.Compile(uriPattern)
	if err != nil {
		logger.Error("编译正则表达式[%s]异常,%s", uriPattern, err)
		return nil
	}
	return UriPatternMatcher(matcher)
}

func (this *regexUriPatternMatcher) match(uri string) bool {
	return uri != "" && this.pattern.MatchString(uri)
}
