package collectors

import (
	"context"
	"io"

	"github.com/cirocosta/devents/lib/events"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"

	dockerevents "github.com/docker/docker/api/types/events"
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

func (d Docker) Collect() (<-chan events.ContainerEvent, <-chan error) {
	var cevents = make(chan events.ContainerEvent)
	var cerrors = make(chan error)

	messages, errs := d.docker.Events(context.Background(), types.EventsOptions{})
	go func() {
		defer close(cevents)
		defer close(cerrors)

		for {
			select {
			case err := <-errs:
				if err != nil && err != io.EOF {
					err = errors.Wrapf(err,
						"Unexpected error happened while receiving events")
					cerrors <- err
				}
				return
			case e := <-messages:
				switch e.Type {
				case dockerevents.ContainerEventType:
					cevents <- events.ContainerEvent{
						Action:      e.Action,
						Image:       e.From,
						ContainerId: e.ID,
						TimeNano:    e.TimeNano,
					}
				}
			}
		}
	}()

	return cevents, cerrors
}
