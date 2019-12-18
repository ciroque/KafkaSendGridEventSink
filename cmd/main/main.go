package main

import (
	config2 "KafkaSendGridEventSink/internal/config"
	daemon2 "KafkaSendGridEventSink/internal/daemon"
	"KafkaSendGridEventSink/pkg/eventing"
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
