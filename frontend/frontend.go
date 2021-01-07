package frontend

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/floj/loxy/backend"
	"github.com/floj/loxy/config"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	mtNumRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "loxy",
		Subsystem: "frontends",
		Name:      "num_requests",
		Help:      "The total number of requests processed",
	}, []string{"name"})

	mtDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loxy",
		Subsystem: "frontends",
		Name:      "duration",
		Help:      "The duration that requests took to be processed",
	}, []string{"name"})
)

type Frontends map[string]*Frontend

type Frontend struct {
	Name   string
	Bind   string
	Port   int
	Routes []*Route
	logger *zap.SugaredLogger
}

func (fe *Frontend) Start() (func() error, error) {
	addr := fmt.Sprintf("%s:%d", fe.Bind, fe.Port)
	fe.logger.Infof("Starting frontend %s on %s", fe.Name, addr)

	svr := http.Server{Addr: addr, Handler: fe}
	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() error {
		fe.logger.Infof("Stopping frontend %s", fe.Name)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return svr.Shutdown(ctx)
	}, nil
}

func (fe *Frontend) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mtNumRequests.WithLabelValues(fe.Name).Inc()
	defer func(start int64) {
		mtDuration.WithLabelValues(fe.Name).Observe(float64(time.Now().Unix() - start))
	}(time.Now().Unix())

	for _, r := range fe.Routes {
		if !r.Matches(req, fe.logger) {
			continue
		}
		r.ApplyModifications(req, fe.logger)
		r.Backend.ServeHTTP(w, req)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func NewFrontend(conf config.Frontend, backends backend.Backends, logger *zap.SugaredLogger) (*Frontend, error) {
	fe := Frontend{
		Name:   conf.Name,
		Bind:   conf.Bind,
		Port:   conf.Port,
		logger: logger,
	}

	for _, rc := range conf.Routes {
		logger.Infof("Creating route %s", rc.Name)
		var be http.Handler = backends[rc.Backend]
		if be == nil {
			return nil, fmt.Errorf("Backend '%s' not found", rc.Backend)
		}

		if conf.Middlewares != nil {
			if conf.Middlewares.Logger != nil {
				logger.Infof("Adding middleware logger")
				be = handlers.LoggingHandler(os.Stdout, be)
			}
			if conf.Middlewares.ProxyHeaders != nil {
				logger.Infof("Adding middleware proxy header")
				be = handlers.ProxyHeaders(be)
			}
		}

		r := Route{Name: rc.Name, Backend: be}
		if rc.Condition != nil {
			for _, c := range rc.Condition.Headers {
				m := HeaderMatcher{Field: c.Field}
				m.setEq(append([]string{c.Eq}, c.EqA...)...)
				m.setPrefix(append([]string{c.Prefix}, c.PrefixA...)...)
				m.setSuffix(append([]string{c.Suffix}, c.SuffixA...)...)
				err := m.setRegexp(append([]string{c.Regexp}, c.RegexpA...)...)
				if err != nil {
					return nil, err
				}
				r.Matcher = append(r.Matcher, &m)
			}
			for _, c := range rc.Condition.Paths {
				m := PathMatcher{}
				m.setEq(append([]string{c.Eq}, c.EqA...)...)
				m.setPrefix(append([]string{c.Prefix}, c.PrefixA...)...)
				m.setSuffix(append([]string{c.Suffix}, c.SuffixA...)...)
				err := m.setRegexp(append([]string{c.Regexp}, c.RegexpA...)...)
				if err != nil {
					return nil, err
				}
				r.Matcher = append(r.Matcher, &m)
			}
		}

		if rc.Modification != nil {
			for _, c := range rc.Modification.Headers {
				r.Modifier = append(r.Modifier, &HeaderModifier{c: c})
			}
			for _, c := range rc.Modification.Paths {
				r.Modifier = append(r.Modifier, &PathModifier{c: c})
			}
		}

		fe.Routes = append(fe.Routes, &r)
	}
	return &fe, nil
}
