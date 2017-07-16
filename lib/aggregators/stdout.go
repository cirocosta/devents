package aggregators

import (
	"github.com/cirocosta/devents/lib/events"

	log "github.com/sirupsen/logrus"
)

type Stdout struct{}

func NewStdout() (agg Stdout, err error) {
	return
}

func (s Stdout) Run(evs <-chan events.ContainerEvent, errs <-chan error) {
	for {
		select {
		case err := <-errs:
			log.WithError(err).Info("errored")
		case ev := <-evs:
			log.WithField("event", ev).Info("event received")
		}
	}
}
