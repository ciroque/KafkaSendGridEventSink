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
	AbortChannel chan error
	Settings     *config2.Settings
}

func (server *Server) Run() {
	http.HandleFunc("/healthz/ping", server.handleHealthz)
	http.HandleFunc("/email/event", server.handleEmailEvent)
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

func (server *Server) handleEmailEvent(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		{
			var event eventing.SendGridEvent
			err := json.NewDecoder(request.Body).Decode(&event)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(writer, err.Error())
				return
			}

			fmt.Printf("Got a thing: %#v", event)

			writer.WriteHeader(http.StatusAccepted)
			fmt.Fprint(writer, "Accepted for processing")
		}
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(writer, "Method not allowed")
	}
}
