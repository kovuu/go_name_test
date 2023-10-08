package interfaces

type AgeGeneratorResult struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type GenderGeneratorResult struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type NationalityGeneratorResult struct {
	Count       int       `json:"count"`
	Name        string    `json:"name"`
	Nationality []Country `json:"country"`
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
