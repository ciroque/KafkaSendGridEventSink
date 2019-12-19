package config

import (
	"os"
	"strconv"
	"testing"
)

func TestSettingsLoadFromEnvironment(t *testing.T) {
	expectedPort := "9099"
	expectedHost := "sample-host"
	expectedKafkaBootstrapServers := "bootstrap.servers"
	expectedKafkaTopic := "kafka.topic"

	os.Setenv("PORT", string(expectedPort))
	os.Setenv("HOST", expectedHost)
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", expectedKafkaBootstrapServers)
	os.Setenv("KAFKA_TOPIC", expectedKafkaTopic)

	settings, err := NewSettings()
	if err != nil {
		t.Fatalf("failed to load Settings, %v", err)
	}

	nport, _ := strconv.Atoi(expectedPort)
	if settings.Port != nport {
		t.Fatalf("Failed loading Port. Expected: '%v', Actual: '%v'", expectedPort, settings.Port)
	}

	if settings.Host != expectedHost {
		t.Fatalf("Failed loading Host. Expected: '%v', Actual: '%v'", expectedHost, settings.Host)
	}

	if settings.KafkaBootstrapServers != expectedKafkaBootstrapServers {
		t.Fatalf("Failed loading Port. Expected: '%v', Actual: '%v'", expectedKafkaBootstrapServers, settings.KafkaBootstrapServers)
	}

	if settings.KafkaTopic != expectedKafkaTopic {
		t.Fatalf("Failed loading Port. Expected: '%v', Actual: '%v'", expectedKafkaTopic, settings.KafkaTopic)
	}
}

func TestDefaults(t *testing.T) {
	os.Clearenv()
	settings, err := NewSettings()
	if err != nil {
		t.Fatalf("failed to load Settings, %v", err)
	}

	expectedPort := 80
	if settings.Port != expectedPort {
		t.Fatalf("Failed loading Port. Expected: '%v', Actual: '%v'", expectedPort, settings.Port)
	}

	expectedHost := "0.0.0.0"
	if settings.Host != expectedHost {
		t.Fatalf("Failed loading Host. Expected: '%v', Actual: '%v'", expectedHost, settings.Host)
	}

	expectedKafkaBootstrapServers := "localhost:9092"
	if settings.KafkaBootstrapServers != expectedKafkaBootstrapServers {
		t.Fatalf("Failed loading Kafka Bootstrap Servers. Expected: '%v', Actual: '%v'", expectedKafkaBootstrapServers, settings.KafkaBootstrapServers)
	}

	expectedKafkaTopic := "test"
	if settings.KafkaTopic != expectedKafkaTopic {
		t.Fatalf("Failed loading Kafka Topic. Expected: '%v', Actual: '%v'", expectedKafkaTopic, settings.KafkaTopic)
	}
}
