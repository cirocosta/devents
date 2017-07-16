package lib

import (
	"github.com/cirocosta/devents/lib/aggregators"
	"github.com/cirocosta/devents/lib/collectors"
	"github.com/docker/docker/api/types/events"
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

		switch agg {
		case "fluentd":
			aggregator, err = aggregators.NewFluentd(aggregators.FluentdConfig{
				Host:      cfg.FluentdHost,
				Port:      cfg.FluentdPort,
				TagPrefix: cfg.FluentdTag,
			})
		case "stdout":
			aggregator, err = aggregators.NewStdout()
		case "prometheus":
			aggregator, err = aggregators.NewPrometheus(aggregators.PrometheusConfig{
				Path:   cfg.MetricsPath,
				Port:   cfg.MetricsPort,
				Labels: cfg.MetricsLabel,
			})
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
	var evChannels []chan events.Message
	var errChannels []chan error

	for idx := range dev.aggregators {
		evChannels = append(evChannels, make(chan events.Message, 1))
		defer close(evChannels[idx])
	}

	for idx := range dev.aggregators {
		errChannels = append(errChannels, make(chan error, 1))
		defer close(errChannels[idx])
	}

	for idx, agg := range dev.aggregators {
		go agg.Run(evChannels[idx], errChannels[idx])
	}

	log.Info("starting main ev loop")
	cevents, cerrors := dev.collector.Collect()

	for {
		select {
		case err := <-cerrors:
			log.WithError(err).Error("error received")
			for _, chann := range errChannels {
				chann <- err
			}

			log.WithError(err).Fatal("Errored waiting for events")
			return
		case ev := <-cevents:
			log.Info("event received")
			for idx, chann := range evChannels {
				select {
				case chann <- ev:
					log.Info("event sent to channel", idx)
				default:
					log.Info("didnt send to channel", idx)
				}
			}
		}
	}
}

// Close closes all aggregators and collectors
func (dev Devents) Close() (err error) {
	return
}
