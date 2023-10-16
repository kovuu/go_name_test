package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Env string `env:"ENV"`
	ConfigDatabase
	MigrationsConfig
	KafkaConfig
	RedisConfig
}

//type HTTPServer struct {
//	Address     string        `yaml:"address" env:"HTTP_ADDRESS" env-default:"localhost:8080"`
//	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"4s"`
//	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
//}

type KafkaConfig struct {
	KafkaUrl            string `env:"KAFKA_URL"`
	KafkaPartition      int    `env:"PARTITION"`
	KafkaFIOTopic       string `env:"FIO_TOPIC"`
	KafkaFIOErrorsTopic string `env:"FIO_ERROR_TOPIC"`
}

type ConfigDatabase struct {
	Port     string `env:"DB_PORT" env-default:"5432"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
	User     string `env:"DB_USER" env-default:"user"`
	Password string `env:"DB_PASSWORD"`
}

type MigrationsConfig struct {
	MigrationPath string `env:"MIGRATIONS_PATH"`
	DriverName    string `env:"DB_DRIVER_NAME"`
}

type RedisConfig struct {
	RedisUrl      string `env:"REDIS_URL"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`
}

func MustLoad() *Config {
	var cfg Config

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config")
	}

	return &cfg
}
