package generator_service

import (
	"encoding/json"
	"fmt"
	"go_test/interfaces"
	"io"
	"net/http"
	"os"
)

func GetAgeGeneratorResult(name string, app *interfaces.Application) int {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		app.Logger.Error("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	var ageGeneratorResult = interfaces.AgeGeneratorResult{}
	err = json.Unmarshal(resBody, &ageGeneratorResult)
	if err != nil {
		app.Logger.Error("failed age generator result parsing", err)
	}
	return ageGeneratorResult.Age
}

func GetGenderGeneratorResult(name string, app *interfaces.Application) string {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		app.Logger.Error("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	var genderGeneratorResult = interfaces.GenderGeneratorResult{}
	err = json.Unmarshal(resBody, &genderGeneratorResult)
	if err != nil {
		app.Logger.Error("failed age generator result parsing", err)
	}
	return genderGeneratorResult.Gender
}

func GetNationalityGeneratorResult(name string, app *interfaces.Application) string {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		app.Logger.Error("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	var nationalityGeneratorResult = interfaces.NationalityGeneratorResult{}
	err = json.Unmarshal(resBody, &nationalityGeneratorResult)
	if err != nil {
		app.Logger.Error("failed age generator result parsing", err)
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
