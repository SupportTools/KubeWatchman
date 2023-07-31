package monitoring

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var controllersStarted = false

func SetControllersStarted() {
	controllersStarted = true
}

func StartWebServer(port string, logger *logrus.Logger) {
	logger.Info("Starting web server on port ", port)
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthzHandler)
	r.HandleFunc("/readyz", readyzHandler)
	// Register Prometheus metrics handler
	r.Handle("/metrics", MetricsHandler())

	http.Handle("/", r)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("Failed to start server: ", err)
	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	if controllersStarted {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
