package controllers_test

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"
	"github.com/supporttools/KubeWatchman/controllers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

func TestPodMonitorControllerStart(t *testing.T) {
	// Create a fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// Create a mock factory
	mockFactory := new(MockClusterConnectionFactory)
	mockFactory.On("NewForConfig", mock.Anything).Return(fakeClientset, nil)

	// Use the mock factory to create a Clientset
	clientset, err := mockFactory.NewForConfig(&rest.Config{})
	if err != nil {
		t.Errorf("Failed to create clientset: %v", err)
	}

	// Create a logger and hook for testing
	logger, _ := test.NewNullLogger()

	// Create the PodMonitorController using the clientset
	controller := controllers.NewPodMonitorController(clientset, logger) // Use clientset here directly

	// Start the controller
	if err := controller.Start(); err != nil {
		t.Errorf("Failed to start controller: %v", err)
	}

	// Additional testing logic needed here
}
