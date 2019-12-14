package processing

import (
	"KafkaSendGridEventSink/pkg/eventing"
	"github.com/Sirupsen/logrus"
)

type Producer struct {
	ProducerChannel chan eventing.SendGridEvent
}

func (producer *Producer) Run() {

	logrus.Info("Producer running...")

	for sendGridEvent := range producer.ProducerChannel {
		logrus.Info("Received an event: %v", sendGridEvent)
		go producer.produce(sendGridEvent)
	}
}

func (producer *Producer) produce(event eventing.SendGridEvent) {
	logrus.Info("Producing an event: %v", event)
}
