package main

import (
	arg "github.com/alexflint/go-arg"
	lib "github.com/cirocosta/devents/lib"
	log "github.com/sirupsen/logrus"
)

var (
	config = lib.Config{
		DockerHost:  "unix://var/run/docker.sock",
		FluentdTag:  "devents",
		FluentdHost: "localhost",
		FluentdPort: 24224,
		Aggregator:  []string{},
		MetricsPath: "/metrics",
		MetricsPort: 9103,
	}
)

func main() {
	arg.MustParse(&config)
	var logger = log.WithFields(config.ToLogrusFields())
	if err := config.Validate(); err != nil {
		logger.
			WithError(err).
			Fatal("Invalid configuration. See `devents -h`")
	}

	dev, err := lib.New(config)
	if err != nil {
		logger.
			WithError(err).
			Fatal("Couldn't initialize Devents")
	}
	defer dev.Close()

	logger.Info("starting")
	dev.Run()
}
