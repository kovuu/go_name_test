package consumer

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strconv"
)

var fioTopic string
var kafkaAddress string
var partition int

func init() {
	err := godotenv.Load(".env")

	var exist bool
	fioTopic, exist = os.LookupEnv("FIO_TOPIC")
	if !exist {
		log.Fatalln("Cannot find fio topic")
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

func ConsumeMessages(c chan string) {
	fmt.Println("consume start")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaAddress},
		Topic:     fioTopic,
		Partition: partition,
		MaxBytes:  10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		c <- string(m.Value)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
