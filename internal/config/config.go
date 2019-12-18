package config

import (
	"fmt"
	"os"
	"strconv"
)

type Settings struct {
	Host                  string
	Port                  int
	KafkaBootstrapServers string
	KafkaTopic            string
}

func NewSettings() (*Settings, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	nport, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("unable to parse PORT: %v", err)
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	kafkaBootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if host == "" {
		host = "localhost:9092"
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if host == "" {
		host = "test"
	}

	config := &Settings{host, nport, kafkaBootstrapServers, kafkaTopic}

	return config, nil
}
