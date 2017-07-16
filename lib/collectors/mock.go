package collectors

import (
	"github.com/cirocosta/devents/lib/events"
)

type Mock struct{}

func NewMock() (agg Mock, err error) {
	return
}

func (Mock) Collect() (<-chan events.ContainerEvent, <-chan error) {
	return nil, nil
}
