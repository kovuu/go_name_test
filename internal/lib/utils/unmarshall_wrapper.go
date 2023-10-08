package utils

import (
	"encoding/json"
	"go_test/interfaces"
	"os"
)

func UnmarshallWrapper(personJSON []byte) (interfaces.Person, error) {
	var person interfaces.Person

	err := json.Unmarshal(personJSON, &person)
	if err != nil {
		return interfaces.Person{}, err
	}

	if len(person.Surname) == 0 {
		return person, os.ErrInvalid
	}
	if len(person.Name) == 0 {
		return person, os.ErrInvalid
	}
	return person, nil
}
