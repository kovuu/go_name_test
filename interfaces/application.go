package interfaces

import (
	"go_test/internal/config"
	"log/slog"
)

type PersonProcessingApp struct {
	DB                DataBaseInterface
	Cfg               *config.Config
	Logger            *slog.Logger
	FioConsumer       FioConsumerInterface
	FioFailedProducer FioFailedProducerInterface
	GeneratorService  PersonInfoGenerator
}

type PersonProducerApp struct {
	Cfg         *config.Config
	Logger      *slog.Logger
	FioProducer AddFioProducerInterface
}
