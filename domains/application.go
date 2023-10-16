package domains

import (
	"go_test/interfaces"
	"go_test/internal/config"
	"log/slog"
)

type PersonProcessingApp struct {
	DB                interfaces.DataBaseInterface
	Cfg               *config.Config
	Logger            *slog.Logger
	FioConsumer       interfaces.FioConsumerInterface
	FioFailedProducer interfaces.FioFailedProducerInterface
	GeneratorService  interfaces.PersonInfoGenerator
	PersonHTTPHandler interfaces.PersonHTTPHandlerInterface
	RedisDB           interfaces.RedisClientService
}

type PersonProducerApp struct {
	Cfg         *config.Config
	Logger      *slog.Logger
	FioProducer interfaces.AddFioProducerInterface
}
