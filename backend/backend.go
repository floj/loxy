package backend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/floj/loxy/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	mtNumRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "loxy",
		Subsystem: "backends",
		Name:      "num_requests",
		Help:      "The total number of requests processed",
	}, []string{"name"})

	mtDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loxy",
		Subsystem: "backends",
		Name:      "duration",
		Help:      "The duration that requests took to be processed",
	}, []string{"name"})
)

type Backends map[string]*Backend

type Backend struct {
	Name    string
	logger  *zap.SugaredLogger
	handler http.Handler
}

func NewBackend(c config.Backend, logger *zap.SugaredLogger) (*Backend, error) {
	be := &Backend{Name: c.Name, logger: logger}

	if c.ReverseProxy != nil {
		rp, err := newReverseProxy(*c.ReverseProxy)
		if err != nil {
			return nil, err
		}
		be.handler = rp
		return be, nil
	}

	if c.FileServer != nil {
		be.handler = newFileServer(c.FileServer)
		return be, nil
	}

	if c.Prometheus != nil {
		be.handler = promhttp.Handler()
		return be, nil
	}

	return nil, fmt.Errorf("Only reverse proxy supported as of now")
}

func (be *Backend) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mtNumRequests.WithLabelValues(be.Name).Inc()
	defer func(start int64) {
		mtDuration.WithLabelValues(be.Name).Observe(float64(time.Now().Unix() - start))
	}(time.Now().Unix())

	be.handler.ServeHTTP(w, req)
}
