package utils

import (
	"encoding/json"
	"go_test/types"
	"os"
)

func UnmarshallWrapper(personJSON []byte) (types.Person, error) {
	var person types.Person

	err := json.Unmarshal(personJSON, &person)
	if err != nil {
		return types.Person{}, err
	}

	if len(person.Surname) == 0 {
		return person, os.ErrInvalid
	}
	if len(person.Name) == 0 {
		return person, os.ErrInvalid
	}
	return person, nil
}
