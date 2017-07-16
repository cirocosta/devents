package aggregators

import (
	"strconv"

	"github.com/docker/docker/api/types/events"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type FluentdConfig struct {
	Host      string
	Port      int
	TagPrefix string
}

type Fluentd struct {
	fluent    *fluent.Fluent
	logger    *log.Entry
	tagPrefix string
}

func NewFluentd(cfg FluentdConfig) (agg Fluentd, err error) {
	f, err := fluent.New(fluent.Config{
		FluentHost: cfg.Host,
		FluentPort: cfg.Port,
		TagPrefix:  cfg.TagPrefix,
	})

	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't instantiate fluentd with config %+v", cfg)
		return
	}

	agg.logger = log.WithField("aggregator", "fluentd")
	agg.fluent = f
	agg.tagPrefix = cfg.TagPrefix

	agg.logger.Info("aggregator initialized")
	return
}

func ConvertEventToMap(ev events.Message) map[string]string {
	var evMap = map[string]string{
		"id":       ev.ID,
		"actorId":  ev.Actor.ID,
		"status":   ev.Status,
		"from":     ev.From,
		"type":     ev.Type,
		"action":   ev.Action,
		"timeNano": strconv.FormatInt(ev.TimeNano, 10),
	}

	for k, v := range ev.Actor.Attributes {
		evMap["attrs."+k] = v
	}

	return evMap
}

func (f Fluentd) Run(evs <-chan events.Message, errs <-chan error) {
	var prefix = f.tagPrefix + ".container"

	f.logger.Info("listening to events")
	for {
		select {
		case err := <-errs:
			f.logger.WithError(err).Info("errored")
		case ev := <-evs:
			f.logger.Info("received evt")
			err := f.fluent.Post(prefix, ConvertEventToMap(ev))
			if err != nil {
				f.logger.
					WithError(err).
					Error("Errored sending ev data to fluentd")
			}
			f.logger.Info("evt sent to fluentd")
		}
	}

	return
}
