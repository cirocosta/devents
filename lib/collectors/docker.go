package collectors

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/cirocosta/devents/lib/events"

	dockerevents "github.com/docker/docker/api/types/events"
	log "github.com/sirupsen/logrus"
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

func (d Docker) Collect() (<-chan events.ContainerEvent,  <-chan error) {
	messages, errs := d.docker.Events(context.Background(), types.EventsOptions{})
	for {
		select {
		case err := <-errs:
			log.WithError(err).Error("Errored.")
		case e := <-messages:
			switch e.Type {
			case dockerevents.ContainerEventType:
			}
			spew.Dump(e)
		}
	}

	return nil, nil
}

