package fio_consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go_test/interfaces"
	"go_test/internal/config"
	"go_test/internal/lib/utils"
	"log"
)

type FioConsumer struct {
	Reader *kafka.Reader
}

func New(cfg *config.Config) *FioConsumer {
	fmt.Println("config", cfg.KafkaFIOTopic)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{cfg.KafkaUrl},
		Topic:     cfg.KafkaFIOTopic,
		Partition: cfg.KafkaPartition,
		MaxBytes:  10e6,
	})
	return &FioConsumer{Reader: reader}
}

func (consumer *FioConsumer) Process(app *interfaces.Application) error {
	for {
		m, err := consumer.Reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		app.Logger.Info("New Value", m.Value)
		//c <- string(m.Value)

		message := m.Value
		person, err := utils.UnmarshallWrapper(message)
		if err != nil {
			log.Printf("Cannot parse a person %s", person)
			personFailedJSON := utils.CreatePersonErrorJSON(person)
			_ = personFailedJSON
			fmt.Println("app", app)
			app.FioFailedProducer.Process(personFailedJSON, app)
			//fmt.Println(personFailedJSON)
			//app.Process(personFailedJSON)
		} else {
			fmt.Println(person)
			app.DB.SavePerson(person, app)
			//db.SavePersonToDB(person)
		}
	}
	return nil
}

func (consumer *FioConsumer) Close() {
	if err := consumer.Reader.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
