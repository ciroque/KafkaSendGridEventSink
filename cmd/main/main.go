package main

import (
	config2 "KafkaSendGridEventSink/internal/config"
	"KafkaSendGridEventSink/internal/event"
	"KafkaSendGridEventSink/internal/web"
	"KafkaSendGridEventSink/pkg/eventing"
	"github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type Main struct {
	AbortChannel    chan error
	ProducerChannel chan []eventing.SendGridEvent
	Settings        *config2.Settings
}

func main() {
	settings, err := config2.NewSettings()
	if err != nil {
		logrus.Fatal("error starting web server", err)
	}

	main := Main{
		AbortChannel:    make(chan error),
		ProducerChannel: make(chan []eventing.SendGridEvent),
		Settings:        settings,
	}

	defer close(main.AbortChannel)
	defer close(main.ProducerChannel)

	main.Run()
}

func (main *Main) startProcessingServer() {
	producer := event.Writer{
		AbortChannel: main.AbortChannel,
		EventChannel: main.ProducerChannel,
		Settings:     main.Settings,
	}

	go producer.Run()
}

func (main *Main) startWebServer() {
	server := web.Server{
		AbortChannel:    main.AbortChannel,
		ProducerChannel: main.ProducerChannel,
		Settings:        main.Settings,
	}

	go server.Run()
}

func (main *Main) Run() {
	main.startProcessingServer()
	main.startWebServer()

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)

	logrus.Info("Running, awaiting signal.")

	select {
	case <-sigTerm:
		{
			logrus.Info("Exiting per SIGTERM")
		}
	case err := <-main.AbortChannel:
		{
			logrus.Error(err)
		}
	}
}
