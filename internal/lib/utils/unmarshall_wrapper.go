package utils

import (
	"encoding/json"
	"go_test/models"
	"os"
)

func UnmarshallWrapper(personJSON []byte) (models.Person, error) {
	var person models.Person

	err := json.Unmarshal(personJSON, &person)
	if err != nil {
		return models.Person{}, err
	}

	if len(person.Surname) == 0 {
		return person, os.ErrInvalid
	}
	if len(person.Name) == 0 {
		return person, os.ErrInvalid
	}
	return person, nil
}
