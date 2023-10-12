package interfaces

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"go_test/domains"
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
	SavePerson(person domains.Person) (int64, error)
	GetPersons(params map[string]string) ([]domains.Person, error)
	GetPersonByID(id int64) (*domains.Person, error)
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
}
