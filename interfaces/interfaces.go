package interfaces

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
	Process(app *PersonProcessingApp) error
	Close()
}

type FioFailedProducerInterface interface {
	Process(personFailedJSON []byte, app *PersonProcessingApp)
	Close()
}

type DataBaseInterface interface {
	SavePerson(person Person, app *PersonProcessingApp) (int64, error)
}

type PersonInfoGenerator interface {
	GetAgeGeneratorResult(name string) int
	GetGenderGeneratorResult(name string) string
	GetNationalityGeneratorResult(name string) string
}

type AddFioProducerInterface interface {
	Process(app *PersonProducerApp)
	Close()
}
