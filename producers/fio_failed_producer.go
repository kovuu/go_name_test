package producers

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strconv"
)

var topicFIOERROR string
var partition int
var kafkaAddress string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no env file", err)
	}
	var exist bool
	topicFIOERROR, exist = os.LookupEnv("FIO_ERROR_TOPIC")
	if !exist {
		log.Fatalln("Cannot find topicFIOERROR")
	}
	kafkaAddress, exist = os.LookupEnv("KAFKA_URL")
	if !exist {
		log.Fatalln("Cannot find apacheAdress")
	}

	p, exist := os.LookupEnv("PARTITION")
	if !exist {
		log.Fatalln("Cannot find partitions")
	}
	buf, err := strconv.Atoi(p)
	if err != nil {
		log.Fatalln("Cannot parse partitions")
		return
	} else {
		partition = buf
	}
}

func ProduceFIOError(personFailedJSON []byte) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, topicFIOERROR, partition)
	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}

	writeMessage(conn, personFailedJSON)

	if err != nil {
		log.Fatal("failed to write messages: ", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func writeMessage(conn *kafka.Conn, person []byte) {
	_, err := conn.WriteMessages(
		kafka.Message{Value: person},
	)

	if err != nil {
		log.Fatal("Ошибка записи", err)
	}
}
