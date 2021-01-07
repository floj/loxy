package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/floj/loxy/backend"
	"github.com/floj/loxy/config"
	"github.com/floj/loxy/frontend"
	"go.uber.org/zap"
)

func main() {
	conf := flag.String("config", "config.hcl", "location of the config file")
	flag.Parse()

	zl, _ := zap.NewDevelopment()
	defer zl.Sync()
	logger := zl.Sugar()

	logger.Infof("Loading config %s", *conf)
	c, err := config.Load(*conf)
	if err != nil {
		logger.Fatalf("Could not load config: %v", err)
	}
	logger.Infof("Config: %+v", c)
	err = run(c, logger)
	if err != nil {
		logger.Fatal(err)
	}

}

func run(c *config.Config, logger *zap.SugaredLogger) error {

	bes := make(backend.Backends)
	for _, b := range c.Backends {
		logger.Infof("Creating backend %s", b.Name)
		be, err := backend.NewBackend(b, logger)
		if err != nil {
			return fmt.Errorf("Could not create backed %s: %w", b.Name, err)
		}
		bes[be.Name] = be
	}

	fes := make(frontend.Frontends)
	for _, f := range c.Frontends {
		logger.Infof("Creating frontend %s", f.Name)
		fe, err := frontend.NewFrontend(f, bes, logger)
		if err != nil {
			return fmt.Errorf("Could not create frontend %s: %w", f.Name, err)
		}
		fes[f.Name] = fe
	}

	var stopFns []func() error
	for _, fe := range fes {
		stopFn, err := fe.Start()
		if err != nil {
			return fmt.Errorf("Could not start frontend %s: %w", fe.Name, err)
		}
		stopFns = append(stopFns, stopFn)
	}

	quitC := make(chan os.Signal, 1)
	signal.Notify(quitC, syscall.SIGTERM, syscall.SIGINT)
	logger.Infof("Running")
	<-quitC

	for _, fn := range stopFns {
		err := fn()
		if err != nil {
			return err
		}
	}
	logger.Infof("Stopped")
	return nil
}
