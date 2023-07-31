package controllers_test

import (
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

type MockClusterConnectionFactory struct {
	mock.Mock
}

func (m *MockClusterConnectionFactory) NewForConfig(config *rest.Config) (*fake.Clientset, error) {
	args := m.Called(config)
	return args.Get(0).(*fake.Clientset), args.Error(1)
}
