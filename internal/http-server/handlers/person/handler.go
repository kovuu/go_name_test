package person

import (
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go_test/domains"
	"go_test/models"
	"strconv"
)

type Handler struct {
	App *domains.PersonProcessingApp
}

func New(app *domains.PersonProcessingApp) *Handler {
	return &Handler{App: app}
}

func (personHandler *Handler) GetPersons(c *routing.Context) error {
	argsMap := parseQueryParamsToMap(c.QueryArgs())
	cachedData, err := personHandler.App.RedisDB.GetFromCache(c.QueryArgs().String())
	if err != nil {
		personHandler.App.Logger.Info("get error while loading data from cache", err)
	}
	if len(cachedData) != 0 {
		c.Response.SetBody(cachedData)
		return nil
	} else {
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
			personHandler.App.RedisDB.SetToCache(c.QueryArgs().String(), personsJSON)
			c.Response.SetBody(personsJSON)
		}
	}

	return nil
}

func (personHandler *Handler) GetPersonByID(c *routing.Context) error {
	if len(c.Param("id")) != 0 {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			personHandler.App.Logger.Info("Cannot parse id param")
			c.Response.SetStatusCode(500)
			c.Response.SetBodyString(err.Error())
			return err
		}
		cacheResult, err := personHandler.App.RedisDB.GetFromCache(c.Param("id"))
		if err != nil {
			personHandler.App.Logger.Info("no data in cache")
		}

		if len(cacheResult) != 0 {
			c.Response.SetBody(cacheResult)
			return nil
		} else {
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
			personHandler.App.RedisDB.SetToCache(string(rune(person.Id)), response)
			c.Response.SetBody(response)
			return nil
		}
	}
	return nil
}

func (personHandler *Handler) SavePerson(c *routing.Context) error {
	var person models.Person
	err := json.Unmarshal(c.Request.Body(), &person)
	if err != nil {
		personHandler.App.Logger.Info("Cannot parse person")
		c.Response.SetStatusCode(500)
		c.Response.SetBodyString(err.Error())
		return err
	}

	if person.Age == 0 {
		person.Age = personHandler.App.GeneratorService.GetAgeGeneratorResult(person.Name)
	}
	if len(person.Gender) == 0 {
		person.Gender = personHandler.App.GeneratorService.GetGenderGeneratorResult(person.Name)
	}
	if len(person.Nationality) == 0 {
		person.Nationality = personHandler.App.GeneratorService.GetNationalityGeneratorResult(person.Name)
	}

	_, err = personHandler.App.DB.SavePerson(person)
	if err != nil {
		personHandler.App.Logger.Info("Cannot save person to database")
		c.Response.SetStatusCode(500)
		c.Response.SetBodyString(err.Error())
		return err
	}
	return nil
}

func (personHandler *Handler) UpdatePerson(c *routing.Context) error {
	var person models.Person
	err := json.Unmarshal(c.Request.Body(), &person)
	if err != nil {
		personHandler.App.Logger.Info("Cannot parse person")
		c.Response.SetStatusCode(500)
		c.Response.SetBodyString(err.Error())
		return err
	}
	err = personHandler.App.DB.UpdatePerson(person)
	if err != nil {
		personHandler.App.Logger.Info("cannot update person", err)
		c.Response.SetStatusCode(500)
		c.Response.SetBodyString(err.Error())
		return err
	}
	c.Response.SetBodyString("Person has been updated successfully")
	return nil
}

func (personHandler *Handler) DeletePerson(c *routing.Context) error {
	if len(c.Param("id")) != 0 {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			personHandler.App.Logger.Info("Cannot parse id param")
			c.Response.SetStatusCode(500)
			c.Response.SetBodyString(err.Error())
			return err
		}
		err = personHandler.App.DB.DeletePersonByID(int64(id))
		if err != nil {
			personHandler.App.Logger.Info("Cannot delete person")
			c.Response.SetStatusCode(500)
			c.Response.SetBodyString(err.Error())
			return err
		}
		c.Response.SetBodyString("Person has	 been deleted!")
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
	return args

}
