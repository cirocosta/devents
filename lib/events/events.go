package events

import (
	"time"
)

type ContainerEvent struct {
	Action string
	Image string
	ContainerId string
	TimeNano time.Time
}
