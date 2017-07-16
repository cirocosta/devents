package lib

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	FluentdHost string   `arg:"help:fluentd host to connect to"`
	FluentdTag  string   `arg:"help:fluentd tag to add to the messages"`
	FluentdPort int      `arg:"help:fluentd port to connect to"`
	DockerHost  string   `arg:"env,help:docker daemon to connect to"`
	Aggregator  []string `arg:"-a,separate,help:aggregators to use"`
}

func (a Config) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"fluentd-host": a.FluentdHost,
		"fluentd-tag":  a.FluentdTag,
		"fluentd-port": a.FluentdPort,
		"docker-host":  a.DockerHost,
		"aggregator":   a.Aggregator,
	}
}

func (a Config) Validate() (err error) {
	if len(a.Aggregator) == 0 {
		err = errors.New(
			"At least one aggregator must be specified.")
		return
	}

	if a.DockerHost == "" {
		err = errors.New(
			"A non-empty docker-host must be specified")
		return
	}

	return
}
