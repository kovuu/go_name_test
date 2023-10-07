package db

import (
	"fmt"
	"go_test/types"
)

func SavePersonToDB(person types.Person) {
	fmt.Println(person, "saved!")
}
