package event

import (
	"KafkaSendGridEventSink/internal/config"
	"KafkaSendGridEventSink/pkg/eventing"
	"bytes"
	"github.com/Sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Writer struct {
	AbortChannel chan error
	EventChannel chan []eventing.SendGridEvent
	Producer     *kafka.Producer
	Settings     *config.Settings
}

func (writer *Writer) Run() {
	configMap := kafka.ConfigMap{
		"bootstrap.servers": writer.Settings.KafkaBootstrapServers,
	}

	logrus.Info("the kafka configMap: %v", configMap)

	p, err := kafka.NewProducer(&configMap)
	if err != nil {
		writer.AbortChannel <- err
		return
	}

	logrus.Info(p.String())

	writer.Producer = p

	logrus.Info("Writer running.")

	go writer.deliveryReporter()

	for events := range writer.EventChannel {
		go writer.produce(events)
	}
}

func (writer *Writer) deliveryReporter() {
	logrus.Info("Delivery Reporter running.")
	for event := range writer.Producer.Events() {
		switch ev := event.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error == nil {

			} else {
				logrus.Infof(
					"Message Delivered. Topic(%v) Partition(%v) Offset(%v)",
					*ev.TopicPartition.Topic,
					ev.TopicPartition.Partition,
					ev.TopicPartition.Offset,
				)
			}
		default:
			logrus.Infof("Ignored event: %v", ev)
		}
	}
	logrus.Info("Delivery Reporter terminating.")
}

func (writer *Writer) produce(events []eventing.SendGridEvent) {
	for _, event := range events {
		logrus.Infof("Producing an event: %#v", events)
		var binary bytes.Buffer
		if err := event.Serialize(&binary); err != nil {
			logrus.Errorf("error serializing the message, %v", err)
			return
		}

		logrus.Debugf("the binary: %v : %v", binary, len(binary.Bytes()))

		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &writer.Settings.KafkaTopic, Partition: kafka.PartitionAny},
			Value:          binary.Bytes(),
		}

		writer.Producer.ProduceChannel() <- message
	}
}
