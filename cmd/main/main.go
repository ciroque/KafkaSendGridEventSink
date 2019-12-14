package main

import (
	config2 "KafkaSendGridEventSink/internal/config"
	"KafkaSendGridEventSink/internal/web"
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type Main struct {
	AbortChannel chan error
	Settings     *config2.Settings
}

func main() {
	settings, err := config2.NewSettings()
	if err != nil {
		logrus.Fatal("error starting web server", err)
	}

	main := Main{
		AbortChannel: make(chan error),
		Settings:     settings,
	}

	main.Run()
}

func (main *Main) startWebServer() {
	server := web.Server{
		AbortChannel: main.AbortChannel,
		Settings:     main.Settings,
	}

	go server.Run()
}

func (main *Main) Run() {
	main.startWebServer()

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)

	fmt.Println("Running...")

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
