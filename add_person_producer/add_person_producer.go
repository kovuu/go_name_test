package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strings"
)

const topicFIO = "FIO"
const partition = 0

type Person struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type CloseProgramErr struct {
}

func (err CloseProgramErr) Error() string {
	return "Закрытие программы"
}

func main() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topicFIO, partition)
	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}
	reading := true
	for reading {
		person, err := readUserData()
		if err != nil {
			reading = false
			return
		}
		personJSOn, err := json.Marshal(person)

		if err != nil {
			log.Fatal("failed to form person json", err)
			return
		}
		go writeMessage(conn, personJSOn)
	}

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
		log.Fatal("Ошибка записи")
	}
}

func readUserData() (Person, error) {
	fmt.Println("Введите информацию о человеке")
	fmt.Print("Имя: ")
	name, err := readString()
	if err != nil {
		log.Fatal("Не удалось прочитать имя")
		return Person{}, err

	}
	fmt.Print("Фамилия: ")
	lastName, err := readString()
	if err != nil {
		log.Fatal("Не удалось прочитать фамилию")
		return Person{}, err

	}
	fmt.Print("Отчество: ")
	patronymic, err := readString()
	if err != nil {
		log.Fatal("Не удалось прочитать отчество")
		return Person{}, err
	}

	person := Person{Name: name, Surname: lastName, Patronymic: patronymic}
	return person, nil
}

func readString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("failed while reading")
		return "", bufio.ErrInvalidUnreadByte
	}

	if strings.TrimSuffix(input, "\n") == "exit" {
		log.Fatal("Выход из программы")
		return "", CloseProgramErr{}
	}
	input = strings.TrimSuffix(input, "\n")
	return input, nil
}

func init() {

}
