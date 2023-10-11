package main

import (
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go_test/domains"
	"go_test/internal/config"
	"go_test/internal/http-server/handlers/person"
	"go_test/internal/services/fio_consumer"
	"go_test/internal/services/generator_service"
	"go_test/internal/services/kafka_fio_errors_producer"
	"go_test/internal/storage/postgres"
	"log/slog"
	"os"
)

type Query struct {
	Filter string `json:"filter"`
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	app := &domains.PersonProcessingApp{}
	app.Cfg = config.MustLoad()

	app.Logger = setupLogger(app.Cfg.Env)
	app.Logger.Debug("debug messages are enabled")

	app.FioConsumer = fio_consumer.New(app)
	app.FioFailedProducer = kafka_fio_errors_producer.New(app)
	var err error
	app.DB, err = postgres.New(app)
	if err != nil {
		app.Logger.Error("db init error", err)
	}
	app.PersonHTTPHandler = person.New(app)
	app.GeneratorService = generator_service.New(app)
	router := routing.New()
	router.Get("/", func(c *routing.Context) error {
		fmt.Fprintf(c, "Hello world")
		return nil
	})
	router.Get("/persons", app.PersonHTTPHandler.GetPersons)

	router.Get("/persons/<id>", app.PersonHTTPHandler.GetPersonByID)
	go fasthttp.ListenAndServe(":8080", router.HandleRequest)

	err = app.FioConsumer.Process()

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
