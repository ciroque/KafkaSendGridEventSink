package daemon

import (
	config2 "kafka-sendgrid-event-sink/internal/config"
	"kafka-sendgrid-event-sink/internal/event"
	"kafka-sendgrid-event-sink/internal/web"
	"kafka-sendgrid-event-sink/pkg/eventing"
	"github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type Daemon struct {
	AbortChannel    chan error
	ProducerChannel chan []eventing.SendGridEvent
	Settings        *config2.Settings
}

func (daemon *Daemon) Run() {
	daemon.startProcessingServer()
	daemon.startWebServer()

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)

	logrus.Info("Running, awaiting signal.")

	select {
	case <-sigTerm:
		{
			close(daemon.ProducerChannel)
			logrus.Info("Exiting per SIGTERM")
		}
	case err := <-daemon.AbortChannel:
		{
			logrus.Error(err)
		}
	}
}

func (daemon *Daemon) startProcessingServer() {
	producer := event.Writer{
		AbortChannel:    daemon.AbortChannel,
		ProducerChannel: daemon.ProducerChannel,
		Settings:        daemon.Settings,
	}

	go producer.Run()
}

func (daemon *Daemon) startWebServer() {
	server := web.Server{
		AbortChannel:    daemon.AbortChannel,
		ProducerChannel: daemon.ProducerChannel,
		Settings:        daemon.Settings,
	}

	go server.Run()
}
