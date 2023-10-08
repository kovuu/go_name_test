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
	Process(app *Application) error
	Close()
}

type FioFailedProducerInterface interface {
	Process(personFailedJSON []byte, app *Application)
	Close()
}

type DataBaseInterface interface {
	SavePerson(person Person, app *Application) (int64, error)
}
