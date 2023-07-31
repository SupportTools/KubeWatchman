package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	controllersUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "controllers_up",
			Help: "Status of the controllers",
		},
		[]string{"controller"},
	)
)

func init() {
	prometheus.MustRegister(controllersUp)
}

func ControllerStatus(controller string, status bool) {
	if status {
		controllersUp.WithLabelValues(controller).Set(1)
	} else {
		controllersUp.WithLabelValues(controller).Set(0)
	}
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
