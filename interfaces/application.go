package interfaces

import (
	"go_test/internal/config"
	"log/slog"
)

type Application struct {
	DB                DataBaseInterface
	Cfg               *config.Config
	Logger            *slog.Logger
	FioConsumer       FioConsumerInterface
	FioFailedProducer FioFailedProducerInterface
	GeneratorService  PersonInfoGenerator
}

func (app *Application) Process(st []byte) {
	//switch st {
	//case "fio_consumer":
	//	err := app.FioConsumer.Process(app)
	//	if err != nil {
	//case "fio_error_producer":
	app.FioFailedProducer.Process(st, app)
	//}
}

//}
