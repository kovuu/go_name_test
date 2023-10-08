package main

import (
	"go_test/interfaces"
	"go_test/internal/config"
	"go_test/internal/services/fio_consumer"
	"go_test/internal/services/kafka_fio_errors_producer"
	"go_test/internal/storage/postgres"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	app := &interfaces.Application{}
	app.Cfg = config.MustLoad()

	app.Logger = setupLogger(app.Cfg.Env)
	app.Logger.Debug("debug messages are enabled")

	app.FioConsumer = fio_consumer.New(app.Cfg)
	app.FioFailedProducer = kafka_fio_errors_producer.New(app.Cfg)
	var err error
	app.DB, err = postgres.New(app.Cfg)
	if err != nil {
		app.Logger.Error("db init error", err)
	}

	err = app.FioConsumer.Process(app)
	if err != nil {
		app.Logger.Error("kafka consumer error", err)
	}

	//storage, err := postgres.New(app.Cfg)
	//if err != nil {
	//	app.Logger.Error("failed to init storage", sl.Err(err))
	//	os.Exit(1)
	//}
	//
	//app.DB = storage

	//incomeMessages := make(chan string)
	//go fio_consumer.ConsumeMessages(incomeMessages)
	//
	//func() {
	//	for {
	//		message := <-incomeMessages
	//		person, err := utils2.UnmarshallWrapper([]byte(message))
	//		if err != nil {
	//			log.Printf("Cannot parse a person %s", person)
	//			personFailedJSON := utils2.CreatePersonErrorJSON(person)
	//			kafka_fio_errors_producer.ProduceFIOError(personFailedJSON)
	//		} else {
	//			fmt.Println(person)
	//			//db.SavePersonToDB(person)
	//		}
	//	}
	//}()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
