package person

import (
	"encoding/json"
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go_test/domains"
	"strconv"
)

type GetPersonHandler struct {
	App *domains.PersonProcessingApp
}

func New(app *domains.PersonProcessingApp) *GetPersonHandler {
	return &GetPersonHandler{App: app}
}

func (personHandler *GetPersonHandler) GetPersons(c *routing.Context) error {
	argsMap := parseQueryParamsToMap(c.QueryArgs())
	persons, err := personHandler.App.DB.GetPersons(argsMap)
	if err != nil {
		personHandler.App.Logger.Info("Cannot select persons from database", err)
	}
	personsJSON, err := json.Marshal(persons)
	if err != nil {
		personHandler.App.Logger.Info("cannot parse persons json", err)
		c.Response.SetStatusCode(500)
		c.Response.SetBody(personsJSON)
	} else {
		c.Response.SetBody(personsJSON)
	}

	return nil
}

func (personHandler *GetPersonHandler) GetPersonByID(c *routing.Context) error {
	if len(c.Param("id")) != 0 {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			personHandler.App.Logger.Info("Cannot parse id param")
			c.Response.SetStatusCode(500)
			c.Response.SetBodyString(err.Error())
			return err
		}
		person, err := personHandler.App.DB.GetPersonByID(int64(id))
		if err != nil {
			personHandler.App.Logger.Info("cannot took person from database")
			c.Response.SetStatusCode(500)
			c.Response.SetBodyString(err.Error())
		}
		response, err := json.Marshal(person)
		if err != nil {
			c.Response.SetStatusCode(500)
			c.Response.SetBodyString("Person object unmarshall failed")
		}
		c.Response.SetBody(response)
	}
	return nil
}

func parseQueryParamsToMap(queryParams *fasthttp.Args) map[string]string {
	args := make(map[string]string)
	if queryParams.Has("filter") {
		args["filter"] = string(queryParams.Peek("filter"))
	}

	if queryParams.Has("orderBy") {
		args["orderBy"] = string(queryParams.Peek("orderBy"))
	}

	if queryParams.Has("limit") {
		args["limit"] = string(queryParams.Peek("limit"))
	}

	if queryParams.Has("offset") {
		args["offset"] = string(queryParams.Peek("offset"))
	}
	fmt.Println("argsMap", args)
	return args

}
