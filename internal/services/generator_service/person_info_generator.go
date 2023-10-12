package generator_service

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go_test/domains"
	"time"
)

const (
	Age         = "https://api.agify.io/?name=%s"
	Gender      = "https://api.genderize.io/?name=%s"
	Nationality = "https://api.nationalize.io/?name=%s"
)

type PersonInfoGenerator struct {
	App    *domains.PersonProcessingApp
	Client *fasthttp.Client
}

func New(app *domains.PersonProcessingApp) *PersonInfoGenerator {

	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	client := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	return &PersonInfoGenerator{App: app, Client: client}

}

func (serv *PersonInfoGenerator) getHttpGeneratorResultRequest(url string, name string) []byte {
	req := fasthttp.AcquireRequest()

	parameterizedUrl := fmt.Sprintf(url, name)
	req.SetRequestURI(parameterizedUrl)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	err := serv.Client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err == nil {
		serv.App.Logger.Debug("debug response: %s\n", resp.Body())
	} else {
		serv.App.Logger.Info("ERR Connection error: %v\n", err)
	}
	fasthttp.ReleaseResponse(resp)

	return resp.Body()
}

func (serv *PersonInfoGenerator) GetAgeGeneratorResult(name string) int {
	resBody := serv.getHttpGeneratorResultRequest(Age, name)
	var ageGeneratorResult = domains.AgeGeneratorResult{}
	err := json.Unmarshal(resBody, &ageGeneratorResult)
	if err != nil {
		serv.App.Logger.Error("failed age generator result parsing", err)
	}
	return ageGeneratorResult.Age
}

func (serv *PersonInfoGenerator) GetGenderGeneratorResult(name string) string {
	resBody := serv.getHttpGeneratorResultRequest(Gender, name)
	var genderGeneratorResult = domains.GenderGeneratorResult{}
	err := json.Unmarshal(resBody, &genderGeneratorResult)
	if err != nil {
		serv.App.Logger.Error("failed age generator result parsing", err)
	}
	return genderGeneratorResult.Gender
}

func (serv *PersonInfoGenerator) GetNationalityGeneratorResult(name string) string {
	resBody := serv.getHttpGeneratorResultRequest(Nationality, name)
	var nationalityGeneratorResult = domains.NationalityGeneratorResult{}
	err := json.Unmarshal(resBody, &nationalityGeneratorResult)
	if err != nil {
		serv.App.Logger.Error("failed age generator result parsing", err)
	}

	var nationality = &domains.Country{}
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
