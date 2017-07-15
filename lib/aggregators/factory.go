package aggregators

import (
	"github.com/pkg/errors"
)

func New(aggregatorType string, config interface{}) (agg Aggregator, err error) {
	switch aggregatorType {
	case "fluentd":
		agg, err = NewFluentd(config.(FluentdConfig))
	case "stdout":
		agg, err = NewStdout()
	case "mock":
		agg, err = NewMock()
	default:
		err = errors.Errorf(
			"Unknown aggregator type %s", aggregatorType)
		return
	}

	return
}
