package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Debug      bool
	KUBECONFIG string
	Port       string
}

func LoadConfigFromEnv() Config {
	config := Config{
		Debug:      parseEnvBool("DEBUG"),
		KUBECONFIG: getEnvOrDefault("KUBECONFIG", "~/.kube/config"),
		Port:       getEnvOrDefault("PORT", "8080"),
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		log.Printf("Failed to parse environment variable %s: %v. Using default value: %d", key, err, defaultValue)
		return defaultValue
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
