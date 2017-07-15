package lib

import (
	"github.com/cirocosta/devents/lib/aggregators"
	"github.com/cirocosta/devents/lib/collectors"
)

type Devents struct {
	collector   collectors.Collector
	aggregators []aggregators.Aggregator
}

func New(cfg Config) (dev Devents, err error) {
	return
}

// Close closes all aggregators and collectors
func (dev Devents) Close() (err error) {
	return
}
