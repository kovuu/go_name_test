package kafka_fio_errors_producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go_test/domains"
	"log"
)

type FioFailedProducer struct {
	Conn *kafka.Conn
	App  *domains.PersonProcessingApp
}

func New(app *domains.PersonProcessingApp) *FioFailedProducer {
	conn, err := kafka.DialLeader(context.Background(), "tcp", app.Cfg.KafkaUrl, app.Cfg.KafkaFIOErrorsTopic, app.Cfg.KafkaPartition)

	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}
	return &FioFailedProducer{Conn: conn, App: app}
}

func (producer *FioFailedProducer) Process(personFailedJSON []byte) {
	err := producer.writeMessage(personFailedJSON)

	if err != nil {
		producer.App.Logger.Info("failed to write messages: ", err)
	}
}

func (producer *FioFailedProducer) writeMessage(person []byte) error {
	_, err := producer.Conn.WriteMessages(
		kafka.Message{Value: person},
	)

	if err != nil {
		log.Fatal("Ошибка записи", err)
		return err
	}
	return nil
}

func (producer *FioFailedProducer) Close() {
	if err := producer.Conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
