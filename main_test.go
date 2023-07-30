package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supporttools/KubeWatchman/config"
)

func TestLoadConfigFromEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("DEBUG", "true")

	// Expect to load the correct configuration
	expectedCfg := config.Config{
		Debug:      true,
		KUBECONFIG: os.Getenv("KUBECONFIG"),
	}

	// Load configuration
	cfg := config.LoadConfigFromEnv()

	// Assert that the configuration matches what we expected
	assert.Equal(t, expectedCfg, cfg)

	// Unset environment variables
	os.Unsetenv("DEBUG")
}
