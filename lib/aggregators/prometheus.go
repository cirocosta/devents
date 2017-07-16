package aggregators

import (
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types/events"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

type PrometheusConfig struct {
	Path   string
	Port   int
	Labels []string
}

type Prometheus struct {
	labels []string
	port   int
	path   string
	logger *log.Entry

	containerActions *prometheus.CounterVec
}

func NewPrometheus(cfg PrometheusConfig) (agg Prometheus, err error) {
	agg.logger = log.WithField("aggregator", "prometheus")
	agg.port = cfg.Port
	agg.path = cfg.Path
	agg.labels = cfg.Labels

	var containerActionLabels = []string{"action"}
	for _, label := range agg.labels {
		containerActionLabels = append(containerActionLabels, label)
	}

	agg.containerActions = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "container_action",
		Help:      "Docker container actions performed",
		Subsystem: "devents",
	}, []string{"action"})

	prometheus.MustRegister(agg.containerActions)

	agg.logger.Info("aggregator initialized")
	return
}

func (p Prometheus) Run(evs <-chan events.Message, errs <-chan error) {
	var handlerErrChan = make(chan error)

	go func() {
		http.Handle(p.path, promhttp.Handler())
		err := http.ListenAndServe(fmt.Sprintf(":%d", p.port), nil)
		if err != nil {
			handlerErrChan <- err
		}
	}()

	p.logger.Info("listening to events")
	for {
		select {
		case err := <-handlerErrChan:
			p.logger.
				WithError(err).
				Error("metrics HTTP handler failed")
		case err := <-errs:
			p.logger.
				WithError(err).
				Error("events retrieval failed")
		case ev := <-evs:
			switch ev.Type {
			case events.ContainerEventType:
				labelValues := []string{
					ev.Action,
				}

				attrs := ev.Actor.Attributes
				for _, label := range p.labels {
					v, present := attrs[label]
					if !present {
						v = "_none"
					}
					labelValues = append(labelValues, v)
				}
				p.containerActions.WithLabelValues(labelValues...).Inc()
			}

			p.logger.
				WithField("event", ev).
				Info("event received")
		}
	}
}
