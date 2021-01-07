package frontend

import (
	"net/http"

	"go.uber.org/zap"
)

type Route struct {
	Name     string
	Backend  http.Handler
	Matcher  []Matcher
	Modifier []Modifier
}

func (r *Route) ApplyModifications(req *http.Request, logger *zap.SugaredLogger) {
	for _, m := range r.Modifier {
		m.Apply(req, logger)
	}
}

func (r *Route) Matches(req *http.Request, logger *zap.SugaredLogger) bool {
	logger.Debugf("Checking route %s", r.Name)
	for _, m := range r.Matcher {
		logger.Debugf("Checking matcher %+v", m)
		if !m.Matches(req, logger) {
			logger.Debugf("No match")
			return false
		}
	}
	logger.Debugf("Matched")
	return true
}
