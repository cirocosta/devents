package collectors

import (
	_ "github.com/docker/docker/client"
)

type DockerConfig struct{}

type Docker struct{}

func NewDocker(cfg DockerConfig) (agg Docker, err error) {
	return
}
