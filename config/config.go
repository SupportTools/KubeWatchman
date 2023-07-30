package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Debug      bool
	KUBECONFIG string
}

func LoadConfigFromEnv() Config {
	config := Config{
		Debug:      parseEnvBool("DEBUG"),
		KUBECONFIG: os.Getenv("KUBECONFIG"),
	}

	return config
}

func parseEnvInt(key string) int {
	value := os.Getenv(key)
	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		log.Fatalf("Failed to parse environment variable %s: %v", key, err)
	}
	return intValue
}

func parseEnvBool(key string) bool {
	value := os.Getenv(key)
	boolValue := false
	if value == "true" {
		boolValue = true
	}
	return boolValue
}
