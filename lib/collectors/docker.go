package collectors

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type DockerConfig struct{}

type Docker struct {
	docker *client.Client
}

func NewDocker(cfg DockerConfig) (collector Docker, err error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't instantiate docker client")
		return
	}

	collector.docker = cli
	return
}

func (d Docker) Collect() (<-chan events.Message, <-chan error) {
	return d.docker.Events(context.Background(), types.EventsOptions{})
}
