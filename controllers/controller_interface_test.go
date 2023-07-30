package controllers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockController struct {
	started bool
	stopped bool
}

func (m *MockController) Start() error {
	m.started = true
	return nil
}

func (m *MockController) Stop() error {
	m.stopped = true
	return nil
}

func TestControllerStart(t *testing.T) {
	controller := &MockController{}
	err := controller.Start()
	assert.NoError(t, err)
	assert.True(t, controller.started)
}

func TestControllerStop(t *testing.T) {
	controller := &MockController{}
	err := controller.Stop()
	assert.NoError(t, err)
	assert.True(t, controller.stopped)
}
