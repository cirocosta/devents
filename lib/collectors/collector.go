package collectors

import (
	"github.com/cirocosta/devents/lib/events"
)

type Collector interface {
	Collect() (<-chan events.ContainerEvent, <-chan error)
}
