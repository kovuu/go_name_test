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
	"strconv"
	"strings"
)

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

func (storage *Storage) SavePerson(person domains.Person) (int64, error) {
	storage.App.Logger.Info("Data saved", person)

	const op = "storage.postgres.SavePerson"
	query := fmt.Sprintf("INSERT INTO person(name, surname, patronymic, age, gender, nationality) VALUES ('%v', '%v', '%v', '%v', '%v', '%v') RETURNING id",
		person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	fmt.Println("query", query)
	var id int64
	err := storage.DB.QueryRow(query).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return id, nil
}

func (storage *Storage) GetPersons(params map[string]string) ([]domains.Person, error) {
	const op = "storage.postgres.GetPersons"
	query := strings.Builder{}
	query.WriteString("SELECT * FROM person")
	err := parseParamsQueryToSQL(&query, params)
	if err != nil {
		storage.App.Logger.Info("query params parsing error", err)
	}
	fmt.Println("parsedquery", query.String())
	persons := []domains.Person{}

	err = storage.DB.Select(&persons, query.String())
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	fmt.Println("Persons", persons)

	return persons, nil
}

func (storage *Storage) GetPersonByID(id int64) (*domains.Person, error) {
	person := domains.Person{}

	err := storage.DB.Get(&person, "Select * FROM person WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func parseParamsQueryToSQL(initSb *strings.Builder, param map[string]string) error {
	fmt.Println("param", len(param["filter"]))
	if len(param["filter"]) != 0 {
		filterSb := strings.Builder{}
		initSb.WriteString(" WHERE")
		filterString := strings.Split(param["filter"], "=")
		paramName := filterString[0]
		paramValue := filterString[1]
		filterSb.WriteString(fmt.Sprintf(" %s = '%s'", paramName, paramValue))
		initSb.WriteString(filterSb.String())
	}

	dataFlowSB := strings.Builder{}
	if len(param["orderBy"]) != 0 {
		initSb.WriteString(" ORDER BY ")
		filterString := strings.Split(param["orderBy"], "=")
		paramName := filterString[0]
		paramValue := filterString[1]
		dataFlowSB.WriteString(fmt.Sprintf("%s %s", paramName, paramValue))
	} else {
		dataFlowSB.WriteString(" ORDER BY name ASC ")
	}

	if len(param["limit"]) != 0 {
		limitValue, err := strconv.Atoi(param["limit"])
		if err != nil {
			return err
		}
		dataFlowSB.WriteString(fmt.Sprintf(" LIMIT %v", limitValue))
	}

	if len(param["offset"]) != 0 {
		offsetValue, err := strconv.Atoi(param["offset"])
		if err != nil {
			return err
		}
		dataFlowSB.WriteString(fmt.Sprintf(" OFFSET %v", offsetValue))
	}

	initSb.WriteString(dataFlowSB.String())

	fmt.Println("stringRes", initSb.String())
	return nil
}
