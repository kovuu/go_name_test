package main

import (
	"go_test/interfaces"
	"go_test/internal/config"
	"go_test/internal/services/fio_consumer"
	"go_test/internal/services/generator_service"
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
	app := &interfaces.PersonProcessingApp{}
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

	app.GeneratorService = generator_service.New(app)
	err = app.FioConsumer.Process(app)
	if err != nil {
		app.Logger.Error("kafka consumer error", err)
	}
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
