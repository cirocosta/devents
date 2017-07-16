package lib

import (
	"github.com/cirocosta/devents/lib/aggregators"
	"github.com/cirocosta/devents/lib/collectors"
	"github.com/cirocosta/devents/lib/events"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type Devents struct {
	collector   collectors.Collector
	aggregators []aggregators.Aggregator
}

func New(cfg Config) (dev Devents, err error) {
	log.WithField("type", "docker").Info("initializing collector")
	collector, err := collectors.NewDocker(struct{}{})
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't instantiate docker collector")
		return
	}

	for _, agg := range cfg.Aggregator {
		var aggregator aggregators.Aggregator

		log.WithField("type", agg).Info("initializing aggregator")
		switch agg {
		case "fluentd":
			aggregator, err = aggregators.NewFluentd(aggregators.FluentdConfig{
				Host:      cfg.FluentdHost,
				Port:      cfg.FluentdPort,
				TagPrefix: cfg.FluentdTag,
			})
		case "stdout":
			aggregator, err = aggregators.NewStdout()
		default:
			err = errors.Errorf(
				"Unknown aggregator type %s", agg)
			return
		}

		if err != nil {
			err = errors.Wrapf(err,
				"Couldn't instantiate aggregator %s", agg)
			return
		}
		dev.aggregators = append(dev.aggregators, aggregator)
	}

	dev.collector = collector
	return
}

func (dev Devents) Run() {
	var evChannels = make([]chan events.ContainerEvent, len(dev.aggregators))
	var errChannels = make([]chan error, len(dev.aggregators))

	for idx := range dev.aggregators {
		defer close(evChannels[idx])
		defer close(errChannels[idx])
	}

	for idx, aggregator := range dev.aggregators {
		go func() {
			aggregator.Run(evChannels[idx], errChannels[idx])
		}()
	}

	log.Info("starting main ev loop")
	cevents, cerrors := dev.collector.Collect()
	for {
		select {
		case err := <-cerrors:
			for _, chann := range errChannels {
				chann <- err
			}

			log.WithError(err).Fatal("Errored waiting for events")
			return
		case ev := <-cevents:
			for _, chann := range evChannels {
				chann <- ev
			}
			log.Println(ev)
		}
	}
}

// Close closes all aggregators and collectors
func (dev Devents) Close() (err error) {
	return
}
