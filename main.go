package main

import (
	"github.com/cirocosta/devents/lib"

	arg "github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

type CLIArguments struct {
	FluentdHost string `arg:"help:fluentd host to connect to"`
	FluentdTag  string `arg:"help:fluentd tag to add to the messages"`
	FluentdPort int    `arg:"help:fluentd port to connect to"`

	DockerHost string `arg:"env,help:docker daemon to connect to"`
}

func (a CLIArguments) ToLogrusFields() log.Fields {
	return log.Fields{
		"fluentd-host": a.FluentdHost,
		"fluentd-tag":  a.FluentdTag,
		"fluentd-port": a.FluentdPort,
		"docker-host":  a.DockerHost,
	}
}

var (
	args = CLIArguments{
		DockerHost: "unix://var/run/docker.sock",
	}
)

func main() {
	arg.MustParse(&args)
	var logger = log.WithFields(args.ToLogrusFields())

	dev, err := lib.New(lib.Config{})
	if err != nil {
		logger.
			WithError(err).
			Fatal("Couldn't initialize Devents")
	}
	defer dev.Close()

	log.Info("Success! Now waiting for events")
}
