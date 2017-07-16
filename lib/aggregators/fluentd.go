package aggregators

import (
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

type FluentdMessage struct {
	Status   string `json:"status"`
	ID       string `json:"id"`
	ActorId  string `json:"actorId"`
	From     string `json:"from"`
	Type     string `json:"type"`
	Action   string `json:"action"`
	TimeNano int64  `json:"timeNano"`
}

func ConvertEventToFluentdMessage(ev events.Message) FluentdMessage {
	return FluentdMessage{
		ID:       ev.ID,
		ActorId:  ev.Actor.ID,
		Status:   ev.Status,
		From:     ev.From,
		Type:     ev.Type,
		Action:   ev.Action,
		TimeNano: ev.TimeNano,
	}
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
			err := f.fluent.Post(prefix, ConvertEventToFluentdMessage(ev))
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
