package generator_service

import (
	"encoding/json"
	"fmt"
	"go_test/interfaces"
	"io"
	"net/http"
	"os"
)

const (
	Age         = "https://api.agify.io/?name=%s"
	Gender      = "https://api.genderize.io/?name=%s"
	Nationality = "https://api.nationalize.io/?name=%s"
)

type PersonInfoGenerator struct {
	App *interfaces.Application
}

func New(app *interfaces.Application) *PersonInfoGenerator {
	return &PersonInfoGenerator{App: app}
}

func (serv *PersonInfoGenerator) getHttpGeneratorResultRequest(url string, name string) []byte {
	parameterizedUrl := fmt.Sprintf(url, name)
	req, err := http.NewRequest(http.MethodGet, parameterizedUrl, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		serv.App.Logger.Error("client: error making http request: %s\n", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		serv.App.Logger.Error("Cannot read getResponse body", err)
	}
	return resBody
}

func (serv *PersonInfoGenerator) GetAgeGeneratorResult(name string) int {
	resBody := serv.getHttpGeneratorResultRequest(Age, name)
	var ageGeneratorResult = interfaces.AgeGeneratorResult{}
	err := json.Unmarshal(resBody, &ageGeneratorResult)
	if err != nil {
		serv.App.Logger.Error("failed age generator result parsing", err)
	}
	return ageGeneratorResult.Age
}

func (serv *PersonInfoGenerator) GetGenderGeneratorResult(name string) string {
	resBody := serv.getHttpGeneratorResultRequest(Gender, name)
	var genderGeneratorResult = interfaces.GenderGeneratorResult{}
	err := json.Unmarshal(resBody, &genderGeneratorResult)
	if err != nil {
		serv.App.Logger.Error("failed age generator result parsing", err)
	}
	return genderGeneratorResult.Gender
}

func (serv *PersonInfoGenerator) GetNationalityGeneratorResult(name string) string {
	resBody := serv.getHttpGeneratorResultRequest(Nationality, name)
	var nationalityGeneratorResult = interfaces.NationalityGeneratorResult{}
	err := json.Unmarshal(resBody, &nationalityGeneratorResult)
	if err != nil {
		serv.App.Logger.Error("failed age generator result parsing", err)
	}

	var nationality = &interfaces.Country{}
	for _, v := range nationalityGeneratorResult.Nationality {
		if nationality == nil {
			nationality = &v
		} else {
			if nationality.Probability < v.Probability {
				nationality = &v
			}
		}

	}
	return nationality.CountryId
}
