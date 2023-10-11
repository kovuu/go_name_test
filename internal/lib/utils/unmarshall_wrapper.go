package utils

import (
	"encoding/json"
	"go_test/domains"
	"os"
)

func UnmarshallWrapper(personJSON []byte) (domains.Person, error) {
	var person domains.Person

	err := json.Unmarshal(personJSON, &person)
	if err != nil {
		return domains.Person{}, err
	}

	if len(person.Surname) == 0 {
		return person, os.ErrInvalid
	}
	if len(person.Name) == 0 {
		return person, os.ErrInvalid
	}
	return person, nil
}
