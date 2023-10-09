package kafka_fio_errors_producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go_test/interfaces"
	"go_test/internal/config"
	"log"
)

type FioFailedProducer struct {
	Conn *kafka.Conn
}

func New(cfg *config.Config) *FioFailedProducer {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.KafkaUrl, cfg.KafkaFIOErrorsTopic, cfg.KafkaPartition)

	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}
	return &FioFailedProducer{Conn: conn}
}

func (producer *FioFailedProducer) Process(personFailedJSON []byte, app *interfaces.PersonProcessingApp) {
	err := producer.writeMessage(personFailedJSON)

	if err != nil {
		app.Logger.Info("failed to write messages: ", err)
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
