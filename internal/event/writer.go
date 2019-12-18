package event

import (
	"KafkaSendGridEventSink/internal/config"
	"KafkaSendGridEventSink/pkg/eventing"
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Writer struct {
	AbortChannel chan error
	EventChannel chan eventing.SendGridEvent
	Producer     *kafka.Producer
	Settings     *config.Settings
}

func (writer *Writer) Run() {
	bootstrapServers := "192.168.0.2"
	configMap := kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	}

	logrus.Info("the kafka configMap: %v", configMap)

	p, err := kafka.NewProducer(&configMap)
	if err != nil {
		writer.AbortChannel <- err
		return
	}

	logrus.Info(p.String())

	writer.Producer = p

	logrus.Info("Writer running...")

	for sendGridEvent := range writer.EventChannel {
		logrus.Info("Received an event: %v", sendGridEvent)
		go writer.produce(sendGridEvent)
	}
}

func (writer *Writer) produce(event eventing.SendGridEvent) {
	logrus.Info("Producing an event: %#v", event)

	var binary bytes.Buffer
	if err := event.Serialize(&binary); err != nil {
		logrus.Errorf("error getting Codec for schema, %v", err)
		return
	}

	logrus.Infof("the binary: %v : %v", binary, len(binary.Bytes()))

	topic := "test"

	err := writer.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          binary.Bytes(),
	}, nil)
	if err != nil {
		logrus.Error(fmt.Errorf("error writing to topic, %v", err))
	}
}
