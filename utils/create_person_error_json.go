package utils

import (
	"encoding/json"
	"go_test/types"
)

func CreatePersonErrorJSON(person types.Person) []byte {
	personFailed := types.PersonFailed{
		Name:       person.Name,
		Surname:    person.Surname,
		Patronymic: person.Patronymic,
	}
	if len(personFailed.Surname) == 0 {
		personFailed.Error = "Surname field is empty"
	}
	if len(personFailed.Name) == 0 {
		personFailed.Error = "Name field is empty"
	}

	personFailedJson, err := json.Marshal(personFailed)
	if err != nil {
		return nil
	}
	return personFailedJson
}
