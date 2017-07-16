package aggregators

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/pkg/errors"
)

type FluentdConfig struct {
	Host      string
	Port      int
	TagPrefix string
}

type Fluentd struct {
	fluent *fluent.Fluent
}

func NewFluentd(cfg FluentdConfig) (aggregator Fluentd, err error) {
	logger, err := fluent.New(fluent.Config{
		FluentHost: cfg.Host,
		FluentPort: cfg.Port,
		TagPrefix:  cfg.TagPrefix,
	})

	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't instantiate fluentd with config %+v", cfg)
		return
	}

	aggregator.fluent = logger
	return
}
