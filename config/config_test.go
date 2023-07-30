package config

import (
	"os"
	"testing"
)

func TestParseEnvInt(t *testing.T) {
	key := "TEST_INT"
	err := os.Setenv(key, "42")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer os.Unsetenv(key)

	value := parseEnvInt(key)
	if value != 42 {
		t.Fatalf("Expected 42, got %d", value)
	}
}

func TestParseEnvBool(t *testing.T) {
	key := "TEST_BOOL"
	err := os.Setenv(key, "true")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer os.Unsetenv(key)

	value := parseEnvBool(key)
	if value != true {
		t.Fatalf("Expected true, got %v", value)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	debugKey := "DEBUG"
	kubeConfigKey := "KUBECONFIG"
	debugValue := "true"
	kubeConfigValue := "~/.kube/config"

	err := os.Setenv(debugKey, debugValue)
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer os.Unsetenv(debugKey)

	err = os.Setenv(kubeConfigKey, kubeConfigValue)
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer os.Unsetenv(kubeConfigKey)

	config := LoadConfigFromEnv()
	if config.Debug != true {
		t.Fatalf("Expected Debug to be true, got %v", config.Debug)
	}
	if config.KUBECONFIG != kubeConfigValue {
		t.Fatalf("Expected KUBECONFIG to be %s, got %s", kubeConfigValue, config.KUBECONFIG)
	}
}
