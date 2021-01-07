package frontend

import (
	"net/http"
	"strings"

	"github.com/floj/loxy/config"
	"go.uber.org/zap"
)

type Modifier interface {
	Apply(req *http.Request, logger *zap.SugaredLogger)
}

type HeaderModifier struct {
	c config.FieldModifier
}

func (md *HeaderModifier) Apply(req *http.Request, logger *zap.SugaredLogger) {
	if md.c.SetValue != "" {
		logger.Debugf("Setting header %s: %s", md.c.Field, md.c.SetValue)
		req.Header.Set(md.c.Field, md.c.SetValue)
		return
	}

	if md.c.Remove {
		logger.Debugf("Removing header %s", md.c.Field)
		req.Header.Del(md.c.Field)
		return
	}

	if md.c.AddValue != "" {
		logger.Debugf("Adding header %s: %s", md.c.Field, md.c.AddValue)
		req.Header.Add(md.c.Field, md.c.AddValue)
		return
	}
}

type PathModifier struct {
	c config.StringModifier
}

func (md *PathModifier) Apply(req *http.Request, logger *zap.SugaredLogger) {
	for _, v := range md.c.StripPrefix {
		if !strings.HasPrefix(req.URL.Path, v) {
			continue
		}
		logger.Debugf("Stripping prefix %s from %s", v, req.URL.Path)
		req.URL.Path = strings.TrimPrefix(req.URL.Path, v)
		req.RequestURI = strings.TrimPrefix(req.RequestURI, v)
	}

	for _, v := range md.c.StripSuffix {
		if !strings.HasSuffix(req.URL.Path, v) {
			continue
		}
		logger.Debugf("Stripping suffix %s from %s", v, req.URL.Path)
		req.URL.Path = strings.TrimSuffix(req.URL.Path, v)
		req.RequestURI = strings.TrimSuffix(req.RequestURI, v)
	}
}
