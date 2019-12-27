package main

import (
	config2 "kafka-sendgrid-event-sink/internal/config"
	daemon2 "kafka-sendgrid-event-sink/internal/daemon"
	"kafka-sendgrid-event-sink/pkg/eventing"
	"github.com/Sirupsen/logrus"
)

func main() {
	settings, err := config2.NewSettings()
	if err != nil {
		logrus.Fatal("error starting web server", err)
	}

	daemon := daemon2.Daemon{
		AbortChannel:    make(chan error),
		ProducerChannel: make(chan []eventing.SendGridEvent),
		Settings:        settings,
	}

	daemon.Run()
}
