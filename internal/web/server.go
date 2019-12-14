package web

import (
	config2 "KafkaSendGridEventSink/internal/config"
	"fmt"
	"github.com/Sirupsen/logrus"
	"net/http"
)

type Server struct {
	AbortChannel chan error
	Settings *config2.Settings
}

func (server *Server) Run() {
	http.HandleFunc("/healthz/ping", server.handleHealthz)
	address := fmt.Sprintf("%s:%d", server.Settings.Host, server.Settings.Port)
	logrus.Info("Listening on ", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		server.AbortChannel <- fmt.Errorf("error starting http listener", err)
	}
}

func (server *Server) handleHealthz(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "pong")
}
