package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"go_test/domains"
	"go_test/internal/config"
	"go_test/models"
	"strconv"
	"strings"
)

type QueryParams struct {
	FilterParam string
	FilterValue string
	OrderRule   string
	Limit       uint
	Offset      uint
}

type Storage struct {
	DB  *sqlx.DB
	App *domains.PersonProcessingApp
}

func New(app *domains.PersonProcessingApp) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := sqlx.Connect(app.Cfg.DriverName, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		app.Cfg.Host, app.Cfg.Port, app.Cfg.User, app.Cfg.Password, app.Cfg.Name))
	if err != nil {
		fmt.Println("Database connect error", err)
	}

	var sqlDB = db.DB
	err = migrateSQL(sqlDB, app.Cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{DB: db, App: app}, nil
}

func migrateSQL(conn *sql.DB, cfg *config.Config) error {
	driver, _ := pgMigrate.WithInstance(conn, &pgMigrate.Config{})
	fmt.Println("cfg", driver)
	m, err := migrate.NewWithDatabaseInstance(
		cfg.MigrationPath,
		cfg.Name,
		driver,
	)

	if err != nil {
		fmt.Println("migrateErr")
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (storage *Storage) SavePerson(person models.Person) (int64, error) {
	storage.App.Logger.Info("Data saved", person)
	const op = "storage.postgres.SavePerson"
	query, _, err := sqlx.In(`INSERT INTO person(name, surname, patronymic, age, gender, nationality) VALUES (?, ?, ?, ?, ?, ?)`,
		person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	if err != nil {
		storage.App.Logger.Info("Cannot create query", err)
		return 0, err
	}
	query = storage.DB.Rebind(query)
	_, err = storage.DB.Query(query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return 0, nil
}

func (storage *Storage) GetPersons(params map[string]string) ([]models.Person, error) {
	const op = "storage.postgres.GetPersons"
	var queryParams QueryParams

	query, err := parseParamsQueryToSQL(&queryParams, params)

	if err != nil {
		storage.App.Logger.Info("cannot parse param query to sql", err)
		return nil, err
	}
	rows, err := storage.DB.NamedQuery(query, queryParams)
	if err != nil {
		storage.App.Logger.Info("cannot get rows", err)
		return nil, err
	}

	var persons []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.StructScan(&person)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	fmt.Println("persons", persons)

	return persons, nil
}

func (storage *Storage) GetPersonByID(id int64) (*models.Person, error) {
	person := models.Person{}

	err := storage.DB.Get(&person, "Select * FROM person WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (storage *Storage) DeletePersonByID(id int64) error {
	_, err := storage.DB.Exec("DELETE from person WHERE id = $1", id)
	if err != nil {
		storage.App.Logger.Info("cannot delete person", err)
		return err
	}

	return nil
}

func (storage *Storage) UpdatePerson(person models.Person) error {
	storage.App.Logger.Info("Data saved", person)
	const op = "storage.postgres.SavePerson"
	_, err := storage.DB.NamedExec(`UPDATE person set name=:name, surname=:surname, patronymic=:patronymic, age=:age, nationality=:nationality, gender=:gender  WHERE id=:id`, person)
	if err != nil {
		storage.App.Logger.Info("Cannot update person", err)
		return err
	}

	return nil
}

func parseParamsQueryToSQL(paramsObj *QueryParams, getParams map[string]string) (string, error) {
	initSb := strings.Builder{}
	initSb.WriteString("SELECT * FROM person")
	if len(getParams["filter"]) != 0 {
		filterString := strings.Split(getParams["filter"], "=")
		paramName := filterString[0]
		paramValue := filterString[1]
		paramsObj.FilterParam = paramName
		paramsObj.FilterValue = paramValue

		initSb.WriteString(fmt.Sprintf(" WHERE %s=:filtervalue", paramName))
	}

	dataFlowSB := strings.Builder{}
	if len(getParams["orderBy"]) != 0 {
		filterString := strings.Split(getParams["orderBy"], "=")
		paramName := filterString[0]
		paramValue := filterString[1]
		paramsObj.OrderRule = paramName + " " + paramValue
		initSb.WriteString(" ORDER BY :orderrule")

	} else {
		paramsObj.OrderRule = "name by ASC"
		initSb.WriteString(" ORDER BY :orderrule ")
	}

	if len(getParams["limit"]) != 0 {
		limitValue, err := strconv.Atoi(getParams["limit"])
		if err != nil {
			return "", err
		}
		paramsObj.Limit = uint(limitValue)
		initSb.WriteString(" LIMIT :limit")
	}

	if len(getParams["offset"]) != 0 {
		offsetValue, err := strconv.Atoi(getParams["offset"])
		if err != nil {
			return "", err
		}
		paramsObj.Offset = uint(offsetValue)
		initSb.WriteString(" OFFSET :offset")
	}

	initSb.WriteString(dataFlowSB.String())

	fmt.Println("stringRes", initSb.String())
	return initSb.String(), nil
}
