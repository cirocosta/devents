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
		Aggregator: []string{
			"stdout",
		},
	}
)

func main() {
	arg.MustParse(&config)
	var logger = log.WithFields(config.ToLogrusFields())

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
