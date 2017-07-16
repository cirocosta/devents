package collectors

import (
	"github.com/docker/docker/api/types/events"
)

type Collector interface {
	Collect() (<-chan events.Message, <-chan error)
}
