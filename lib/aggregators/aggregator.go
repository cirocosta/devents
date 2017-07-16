package aggregators

import (
	"github.com/cirocosta/devents/lib/events"
)

type Aggregator interface {
	Run(<-chan events.ContainerEvent, <-chan error)
}
