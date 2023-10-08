package fio_consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go_test/interfaces"
	"go_test/internal/config"
	"go_test/internal/lib/utils"
	"go_test/internal/services/generator_service"
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

		message := m.Value
		person, err := utils.UnmarshallWrapper(message)
		if err != nil {
			log.Printf("Cannot parse a person %s", person)
			personFailedJSON := utils.CreatePersonErrorJSON(person)
			_ = personFailedJSON
			app.FioFailedProducer.Process(personFailedJSON, app)
		} else {
			person.Age = generator_service.GetAgeGeneratorResult(person.Name, app)
			person.Gender = generator_service.GetGenderGeneratorResult(person.Name, app)
			person.Nationality = generator_service.GetNationalityGeneratorResult(person.Name, app)
			savedId, err := app.DB.SavePerson(person, app)
			if err != nil {
				fmt.Println("DB ERORR", err)
			} else {
				app.Logger.Info("person has been saved with id ", savedId)
			}
		}
	}
	return nil
}

func (consumer *FioConsumer) Close() {
	if err := consumer.Reader.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
