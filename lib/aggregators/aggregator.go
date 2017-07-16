package aggregators

import (
	"github.com/docker/docker/api/types/events"
)

type Aggregator interface {
	Run(<-chan events.Message, <-chan error)
}
