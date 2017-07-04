package main

import (
	"flag"

	"github.com/cirocosta/devents/lib"

	log "github.com/sirupsen/logrus"
)

var (
	host = flag.String("fluentd-host", "localhost", "Fluentd Host")
	port = flag.Int("fluentd-port", 24224, "Fluentd Port")
	tag  = flag.String("fluentd-tag", "devents.events", "Tag to use for the logs")
)

func main() {
	flag.Parse()
	log.Info("Initializing Devents")

	dev, err := lib.New(lib.Config{
		FluentdHost: *host,
		FluentdPort: *port,
		FluentdTag:  *tag,
	})
	if err != nil {
		log.
			WithError(err).
			Fatal("Couldn't initialize Devents")
	}
	defer dev.Close()

	log.Info("Success! Now waiting for events")
}
