package interfaces

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"go_test/models"
)

type HttpService interface {
	Process()
}

type AgeGeneratorService interface {
	Process()
}

type NationalityGeneratorService interface {
	Process()
}

type SexGeneratorService interface {
	Process()
}

type FioConsumerInterface interface {
	Process() error
	Close()
}

type FioFailedProducerInterface interface {
	Process(personFailedJSON []byte)
	Close()
}

type DataBaseInterface interface {
	SavePerson(person models.Person) (int64, error)
	GetPersons(params map[string]string) ([]models.Person, error)
	GetPersonByID(id int64) (*models.Person, error)
	DeletePersonByID(id int64) error
	UpdatePerson(person models.Person) error
}

type PersonInfoGenerator interface {
	GetAgeGeneratorResult(name string) int
	GetGenderGeneratorResult(name string) string
	GetNationalityGeneratorResult(name string) string
}

type AddFioProducerInterface interface {
	Process()
	Close()
}

type PersonHTTPHandlerInterface interface {
	GetPersons(c *routing.Context) error
	GetPersonByID(c *routing.Context) error
	SavePerson(c *routing.Context) error
	UpdatePerson(c *routing.Context) error
	DeletePerson(c *routing.Context) error
}
