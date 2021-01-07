package frontend

import (
	"net/http"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

type Matcher interface {
	Matches(req *http.Request, logger *zap.SugaredLogger) bool
}

type HeaderMatcher struct {
	StringMatcher
	Field string
}

func (m *HeaderMatcher) Matches(req *http.Request, logger *zap.SugaredLogger) bool {
	values := req.Header.Values(m.Field)
	logger.Debugf("Checking %+v", m)
	for _, v := range values {
		if m.matches(v, logger) {
			logger.Debugf("Match found: %s", v)
			return true
		}
	}
	logger.Debugf("No match found")
	return false
}

type PathMatcher struct {
	StringMatcher
}

func (m *PathMatcher) Matches(req *http.Request, logger *zap.SugaredLogger) bool {
	return m.matches(req.URL.Path, logger)
}

type StringMatcher struct {
	Eq     []string
	Prefix []string
	Suffix []string
	Regexp []*regexp.Regexp
}

func (s *StringMatcher) setEq(v ...string) {
	for _, e := range v {
		if e == "" {
			continue
		}
		s.Eq = append(s.Eq, e)
	}
}

func (s *StringMatcher) setPrefix(v ...string) {
	for _, e := range v {
		if e == "" {
			continue
		}
		s.Prefix = append(s.Prefix, e)
	}
}

func (s *StringMatcher) setSuffix(v ...string) {
	for _, e := range v {
		if e == "" {
			continue
		}
		s.Prefix = append(s.Suffix, e)
	}
}

func (s *StringMatcher) setRegexp(v ...string) error {
	for _, e := range v {
		if e == "" {
			continue
		}
		re, err := regexp.Compile(e)
		if err != nil {
			return err
		}
		s.Regexp = append(s.Regexp, re)
	}
	return nil
}

func (s *StringMatcher) matches(value string, logger *zap.SugaredLogger) bool {
	if len(s.Eq) > 0 {
		matched := false
		for _, t := range s.Eq {
			matched = value == t
			logger.Debugf("eq(%s, %s) => %t", t, value, matched)
			if matched {
				break
			}
		}
		if !matched {
			return false
		}
	}

	if len(s.Prefix) > 0 {
		matched := false
		for _, t := range s.Prefix {
			matched = strings.HasPrefix(value, t)
			logger.Debugf("prefix(%s, %s) => %t", t, value, matched)
			if matched {
				break
			}
		}
		if !matched {
			return false
		}
	}
	if len(s.Suffix) > 0 {
		matched := false
		for _, t := range s.Suffix {
			matched = strings.HasSuffix(value, t)
			logger.Debugf("suffix(%s, %s) => %t", t, value, matched)
			if matched {
				break
			}
		}
		if !matched {
			return false
		}
	}

	if len(s.Regexp) > 0 {
		matched := false
		for _, t := range s.Regexp {
			matched = t.MatchString(value)
			logger.Debugf("regexp(%s, %s) => %t", t, value, matched)
			if matched {
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
