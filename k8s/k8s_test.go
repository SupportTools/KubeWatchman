package k8s_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/supporttools/KubeWatchman/k8s"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type MockCoreV1Interface struct {
	v1.CoreV1Interface
	mock.Mock
}

// real implementation, simply forwards to kubernetes.Clientset
type RealClusterConnection struct {
	*kubernetes.Clientset
}

func (r *RealClusterConnection) RESTClient() rest.Interface {
	return r.Clientset.RESTClient()
}

// mock implementation, can return whatever you want
type MockClusterConnection struct {
	mock.Mock
}

func (m *MockClusterConnection) RESTClient() rest.Interface {
	args := m.Called()
	return args.Get(0).(rest.Interface)
}

func (m *MockClusterConnection) CoreV1() v1.CoreV1Interface {
	args := m.Called()
	return args.Get(0).(v1.CoreV1Interface)
}

type MockClusterConnectionFactory struct {
	mock.Mock
}

func (m *MockClusterConnectionFactory) InClusterConfig() (*rest.Config, error) {
	args := m.Called()
	return args.Get(0).(*rest.Config), args.Error(1)
}

func (m *MockClusterConnectionFactory) NewForConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	args := m.Called(config)
	return args.Get(0).(*kubernetes.Clientset), args.Error(1)
}

func TestCreateClusterConnection(t *testing.T) {
	// Create a new logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	// Create a mock factory
	mockFactory := new(MockClusterConnectionFactory)
	mockFactory.On("InClusterConfig").Return(&rest.Config{}, nil)
	mockFactory.On("NewForConfig", mock.Anything).Return(&kubernetes.Clientset{}, nil)

	// Call CreateClusterConnection with the mock factory
	clientset, err := k8s.CreateClusterConnection(logger, mockFactory) // Adjusted the call here

	// Assert that there was no error
	if err != nil {
		t.Errorf("Failed to create cluster connection: %v", err)
	}

	// Assert that the clientset is not nil
	if clientset == nil {
		t.Errorf("Unexpected nil clientset")
	}
}
