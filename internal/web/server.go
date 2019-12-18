package web

import (
	config2 "KafkaSendGridEventSink/internal/config"
	"KafkaSendGridEventSink/pkg/eventing"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"net/http"
)

type Server struct {
	AbortChannel    chan error
	ProducerChannel chan []eventing.SendGridEvent
	Settings        *config2.Settings
}

func (server *Server) Run() {
	http.HandleFunc("/healthz/ping", server.handleHealthz)
	http.HandleFunc("/email/event", server.handleEmailEvent)
	address := fmt.Sprintf("%s:%d", server.Settings.Host, server.Settings.Port)
	logrus.Info("Listening on ", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		server.AbortChannel <- fmt.Errorf("error starting http listener, %v", err)
	}
}

func (server *Server) handleHealthz(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "pong")
}

func (server *Server) handleEmailEvent(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		{
			var events []eventing.SendGridEvent
			err := json.NewDecoder(request.Body).Decode(&events)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(writer, err.Error())
				return
			}

			logrus.Infof("Received events: %v", events)

			server.ProducerChannel <- events

			writer.WriteHeader(http.StatusAccepted)
			fmt.Fprint(writer, "Accepted")
		}
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(writer, "Method not allowed")
	}
}
