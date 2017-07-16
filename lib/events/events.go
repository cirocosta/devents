package events

type ContainerEvent struct {
	Action      string
	Image       string
	ContainerId string
	TimeNano    int64
}
