package aggregators

import (
	_ "github.com/fluent/fluent-logger-golang/fluent"
)

type FluentdConfig struct{}

type Fluentd struct{}

func NewFluentd(config FluentdConfig) (fluent Fluentd, err error) {
	return
}
