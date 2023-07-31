package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/supporttools/KubeWatchman/config"
	"github.com/supporttools/KubeWatchman/controllers"
	"github.com/supporttools/KubeWatchman/k8s"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type RealClusterConnectionFactory struct{}

func (f *RealClusterConnectionFactory) InClusterConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" || os.Getenv("KUBERNETES_SERVICE_PORT") == "" {
		// If we're not running in a cluster, use the kubeconfig file
		if os.Getenv("KUBECONFIG") != "" {
			return clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		}
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	// Otherwise, use in-cluster configuration
	return rest.InClusterConfig()
}

func (f *RealClusterConnectionFactory) NewForConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}

func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfigFromEnv()

	// Set up logger and log level based on configuration
	logger := logrus.StandardLogger()
	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	factory := &RealClusterConnectionFactory{}
	clientset, err := k8s.CreateClusterConnection(logger, factory)
	if err != nil {
		logrus.Fatalf("Failed to create cluster connection: %v", err)
	}

	err = k8s.CheckClusterConnection(clientset, logger)
	if err != nil {
		logrus.Fatalf("Failed to test cluster connection: %v", err)
	}

	nodeMonitorController := controllers.NewNodeMonitorController(clientset, logger)

	controllersList := []controllers.Controller{
		nodeMonitorController,
	}

	for _, controller := range controllersList {
		logrus.Debug("Starting controller")
		if err := controller.Start(); err != nil {
			logrus.Fatalf("Failed to start controller: %v", err)
		}
	}

	// Channel to catch OS signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal
	<-signalCh

	// Stop controllers when finished
	for _, controller := range controllersList {
		logrus.Debug("Stopping controller")
		//nolint:errcheck
		controller.Stop()
	}

	logrus.Println("Controllers have been stopped gracefully")
}
