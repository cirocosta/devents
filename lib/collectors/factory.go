package collectors

import (
	"github.com/pkg/errors"
)

func New(collectorType string, config interface{}) (collector Collector, err error) {
	switch collectorType {
	case "docker":
		collector, err = NewDocker(config.(DockerConfig))
	default:
		err = errors.Errorf(
			"Unknown collector type %s", collectorType)
		return
	}

	return
}
