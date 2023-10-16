package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go_test/domains"
	"go_test/graph"
	"go_test/internal/config"
	"go_test/internal/http-server/handlers/person"
	"go_test/internal/services/fio_consumer"
	"go_test/internal/services/generator_service"
	"go_test/internal/services/kafka_fio_errors_producer"
	"go_test/internal/storage/postgres"
	"go_test/internal/storage/redis"
	"log"
	"log/slog"
	"net/http"
	"os"
)

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

	app.RedisDB = redis.New(app)
	app.PersonHTTPHandler = person.New(app)
	app.GeneratorService = generator_service.New(app)
	router := routing.New()
	router.Get("/", func(c *routing.Context) error {
		fmt.Fprintf(c, "Hello world")
		return nil
	})
	router.Get("/persons", app.PersonHTTPHandler.GetPersons)
	router.Get("/persons/<id>", app.PersonHTTPHandler.GetPersonByID)
	router.Post("/persons", app.PersonHTTPHandler.SavePerson)
	router.Delete("/persons/<id>", app.PersonHTTPHandler.DeletePerson)
	router.Patch("/persons", app.PersonHTTPHandler.UpdatePerson)
	go fasthttp.ListenAndServe(":8080", router.HandleRequest)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{App: app}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", "8080")
	go http.ListenAndServe(":8081", nil)

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
