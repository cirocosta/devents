package aggregators

import (
	"strconv"

	"github.com/cirocosta/devents/lib/events"
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

func (f Fluentd) Run(evs <-chan events.ContainerEvent, errs <-chan error) {
	var prefix = f.tagPrefix + ".container"

	f.logger.Info("listening to events")
	for {
		select {
		case err := <-errs:
			f.logger.WithError(err).Info("errored")
		case ev := <-evs:
			err := f.fluent.Post(prefix, map[string]string{
				"image":        ev.Image,
				"action":       ev.Action,
				"container_id": ev.ContainerId,
				"time":         strconv.FormatInt(ev.TimeNano, 10),
			})
			if err != nil {
				f.logger.
					WithError(err).
					Error("Errored sending ev data to fluentd")
			}
		}
	}

	return
}
