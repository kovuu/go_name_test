package fio_consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go_test/domains"
	"go_test/internal/lib/utils"
	"log"
)

type FioConsumer struct {
	Reader *kafka.Reader
	App    *domains.PersonProcessingApp
}

func New(app *domains.PersonProcessingApp) *FioConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{app.Cfg.KafkaUrl},
		Topic:     app.Cfg.KafkaFIOTopic,
		Partition: app.Cfg.KafkaPartition,
		MaxBytes:  10e6,
	})
	offset := getLastTopicOffset(app)
	err := reader.SetOffset(offset)
	if err != nil {
		app.Logger.Error("Cannot set consumer offset")
	}

	return &FioConsumer{Reader: reader, App: app}
}

func getLastTopicOffset(app *domains.PersonProcessingApp) int64 {
	conn, err := kafka.DialLeader(context.Background(), "tcp", app.Cfg.KafkaUrl, app.Cfg.KafkaFIOTopic, app.Cfg.KafkaPartition)
	if err != nil {
		app.Logger.Error("Cannot connect to kafka")
	}
	_, last, err := conn.ReadOffsets()
	if err != nil {
		app.Logger.Error("Cannot read offsets")
	}
	err = conn.Close()
	if err != nil {
		app.Logger.Error("Cannot close kafka connection", err)
	}
	return last
}

func (consumer *FioConsumer) Process() error {
	for {
		m, err := consumer.Reader.ReadMessage(context.Background())

		if err != nil {
			consumer.App.Logger.Info("err", err)
			break
		}
		consumer.App.Logger.Info("New Value", m.Value)
		message := m.Value

		person, err := utils.UnmarshallWrapper(message)
		if err != nil {
			consumer.App.Logger.Info("Cannot parse a person %s", person)
			personFailedJSON := utils.CreatePersonErrorJSON(person)
			_ = personFailedJSON
			consumer.App.FioFailedProducer.Process(personFailedJSON)
		} else {
			person.Age = consumer.App.GeneratorService.GetAgeGeneratorResult(person.Name)
			person.Gender = consumer.App.GeneratorService.GetGenderGeneratorResult(person.Name)
			person.Nationality = consumer.App.GeneratorService.GetNationalityGeneratorResult(person.Name)

			if person.Age != 0 && len(person.Gender) != 0 && len(person.Nationality) != 0 {
				savedId, err := consumer.App.DB.SavePerson(person)
				if err != nil {
					consumer.App.Logger.Info("DB ERROR", err)
				} else {
					consumer.App.Logger.Info("person has been saved with id ", savedId)
				}
			} else {
				consumer.App.Logger.Error("Invalid person data")
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
