package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go_test/interfaces"
	"go_test/internal/config"
)

type Storage struct {
	db *sql.DB
}

func New(config *config.Config) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)

	fmt.Println("driver", config.DriverName)
	db, err := sql.Open(config.DriverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = migrateSQL(db, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
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

func (storage *Storage) SavePerson(person interfaces.Person, app *interfaces.Application) {
	app.Logger.Info("Data saved", person)

	//const op = "storage.postgres.SavePerson"
	//query := fmt.Sprintf("INSERT INTO person(url, alias) VALUES ('%v', '%v') RETURNING id", urlToSave, alias)
	//fmt.Println("query", query)
	//var id int64
	//err := s.db.QueryRow(query).Scan(&id)
	//if err != nil {
	//	return 0, fmt.Errorf("%s %w", op, err)
	//}
	//
	//return id, nil
}

//func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
//	const op = "storage.postgres.SaveURL"
//	query := fmt.Sprintf("INSERT INTO url(url, alias) VALUES ('%v', '%v') RETURNING id", urlToSave, alias)
//	fmt.Println("query", query)
//	var id int64
//	err := s.db.QueryRow(query).Scan(&id)
//	if err != nil {
//		return 0, fmt.Errorf("%s %w", op, err)
//	}
//
//	return id, nil
//}
