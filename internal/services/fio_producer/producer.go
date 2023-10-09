package fio_producer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go_test/interfaces"
	"go_test/internal/config"
	"log"
	"os"
	"strings"
)

type CloseProgramErr struct {
}

func (err CloseProgramErr) Error() string {
	return "Закрытие программы"
}

type FioProducer struct {
	Conn *kafka.Conn
}

func New(cfg *config.Config) *FioProducer {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.KafkaUrl, cfg.KafkaFIOErrorsTopic, cfg.KafkaPartition)
	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}
	return &FioProducer{Conn: conn}
}

func (producer *FioProducer) Process(app *interfaces.PersonProducerApp) {
	reading := true
	for reading {
		person, err := producer.readUserData()
		if err != nil {
			reading = false
			return
		}
		personJSOn, err := json.Marshal(person)

		if err != nil {
			log.Fatal("failed to form person json", err)
			return
		}
		go producer.writeMessage(producer.Conn, personJSOn)
	}

}

func (producer *FioProducer) writeMessage(conn *kafka.Conn, person []byte) {
	_, err := conn.WriteMessages(
		kafka.Message{Value: person},
	)

	if err != nil {
		log.Fatal("Ошибка записи")
	}
}

func (producer *FioProducer) readUserData() (interfaces.Person, error) {
	fmt.Println("Введите информацию о человеке")
	fmt.Print("Имя: ")
	name, err := producer.readString()
	if err != nil {
		log.Fatal("Не удалось прочитать имя")
		return interfaces.Person{}, err

	}
	fmt.Print("Фамилия: ")
	lastName, err := producer.readString()
	if err != nil {
		log.Fatal("Не удалось прочитать фамилию")
		return interfaces.Person{}, err

	}
	fmt.Print("Отчество: ")
	patronymic, err := producer.readString()
	if err != nil {
		log.Fatal("Не удалось прочитать отчество")
		return interfaces.Person{}, err
	}

	person := interfaces.Person{Name: name, Surname: lastName, Patronymic: patronymic}
	return person, nil
}

func (producer *FioProducer) readString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("failed while reading")
		return "", bufio.ErrInvalidUnreadByte
	}

	if strings.TrimSuffix(input, "\n") == "exit" {
		log.Fatal("Выход из программы")
		producer.Close()
		return "", CloseProgramErr{}
	}
	input = strings.TrimSuffix(input, "\n")
	return input, nil
}

func (producer *FioProducer) Close() {
	if err := producer.Conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
