package main

import (
	"fmt"
	"go_test/consumer"
	"go_test/db"
	"go_test/producers"
	"go_test/utils"
	"log"
)

func main() {
	incomeMessages := make(chan string)
	go consumer.ConsumeMessages(incomeMessages)

	func() {
		for {
			message := <-incomeMessages
			person, err := utils.UnmarshallWrapper([]byte(message))
			if err != nil {
				log.Printf("Cannot parse a person %s", person)
				personFailedJSON := utils.CreatePersonErrorJSON(person)
				producers.ProduceFIOError(personFailedJSON)
			} else {
				fmt.Println(person)
				db.SavePersonToDB(person)
			}
		}
	}()
}
