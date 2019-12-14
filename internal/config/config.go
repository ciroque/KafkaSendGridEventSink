package config

import (
	"fmt"
	"os"
	"strconv"
)

type Settings struct {
	Host                   string
	Port                   int
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

	config := &Settings{host, nport}

	return config, nil
}
