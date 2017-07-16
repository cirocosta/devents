package aggregators

import (
	"github.com/docker/docker/api/types/events"

	log "github.com/sirupsen/logrus"
)

type Stdout struct {
	logger *log.Entry
}

func NewStdout() (agg Stdout, err error) {
	agg.logger = log.WithField("aggregator", "stdout")
	agg.logger.Info("aggregator initialized")
	return
}

func (s Stdout) Run(evs <-chan events.Message, errs <-chan error) {
	s.logger.Info("listening to events")

	for {
		select {
		case err := <-errs:
			s.logger.WithError(err).Info("errored")
		case ev := <-evs:
			s.logger.WithField("event", ev).Info("event received")
		}
	}
}
