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
	Path string
	Port int
}

type Prometheus struct {
	port   int
	path   string
	logger *log.Entry

	containerActions *prometheus.CounterVec
}

func NewPrometheus(cfg PrometheusConfig) (agg Prometheus, err error) {
	agg.logger = log.WithField("aggregator", "prometheus")
	agg.port = cfg.Port
	agg.path = cfg.Path

	agg.logger.Info("aggregator initialized")
	return
}

func (p Prometheus) MustSetupCounters() {
	p.containerActions = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "devents_container_action",
		Help: "Docker container actions performed",
	}, []string{"action"})

	prometheus.MustRegister(p.containerActions)
}

func (p Prometheus) Run(evs <-chan events.Message, errs <-chan error) {
	var handlerErrChan = make(chan error)

	p.MustSetupCounters()

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
				p.containerActions.WithLabelValues(ev.Action).Inc()
			}

			p.logger.
				WithField("event", ev).
				Info("event received")
		}
	}
}
