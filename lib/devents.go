package lib

import (
	"github.com/docker/docker/client"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/pkg/errors"
)

type Config struct {
	FluentdHost string
	FluentdTag  string
	FluentdPort int
}

type Devents struct {
	docker *client.Client
	writer *fluent.Fluent
}

func New(cfg Config) (dev Devents, err error) {
	if cfg.FluentdHost == "" || cfg.FluentdTag == "" || cfg.FluentdPort == 0 {
		err = errors.Errorf("All configuration must be filled.")
		return
	}

	writer, err := fluent.New(fluent.Config{})
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't instantiate fluent")
		return
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't initialize docker client")
		return
	}

	dev.writer = writer
	dev.docker = cli
	return
}

func (dev Devents) Close() (err error) {
	err = dev.writer.Close()
	return
}
