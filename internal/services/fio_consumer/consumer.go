package fio_consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go_test/interfaces"
	"go_test/internal/lib/utils"
	"log"
)

type FioConsumer struct {
	Reader *kafka.Reader
}

func New(app *interfaces.PersonProcessingApp) *FioConsumer {
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

	return &FioConsumer{Reader: reader}
}

func getLastTopicOffset(app *interfaces.PersonProcessingApp) int64 {
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

func (consumer *FioConsumer) Process(app *interfaces.PersonProcessingApp) error {
	for {
		m, err := consumer.Reader.ReadMessage(context.Background())

		if err != nil {
			app.Logger.Info("err", err)
			break
		}
		app.Logger.Info("New Value", m.Value)
		message := m.Value

		person, err := utils.UnmarshallWrapper(message)
		if err != nil {
			app.Logger.Info("Cannot parse a person %s", person)
			personFailedJSON := utils.CreatePersonErrorJSON(person)
			_ = personFailedJSON
			app.FioFailedProducer.Process(personFailedJSON, app)
		} else {
			person.Age = app.GeneratorService.GetAgeGeneratorResult(person.Name)
			person.Gender = app.GeneratorService.GetGenderGeneratorResult(person.Name)
			person.Nationality = app.GeneratorService.GetNationalityGeneratorResult(person.Name)

			if person.Age != 0 && len(person.Gender) != 0 && len(person.Nationality) != 0 {
				savedId, err := app.DB.SavePerson(person, app)
				if err != nil {
					app.Logger.Info("DB ERROR", err)
				} else {
					app.Logger.Info("person has been saved with id ", savedId)
				}
			} else {
				app.Logger.Error("Invalid person data")
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
