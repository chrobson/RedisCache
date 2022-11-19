package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database interface {
	GetUserById(string) (*Person, error)
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase() (*PostgresDatabase, error) {
	connStr := "user=postgres, dbname=people, password=postgrespw sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDatabase{
		db: db,
	}, nil
}

func (pd *PostgresDatabase) GetUserById(id string) (*Person, error) {
	row := pd.db.QueryRow("SELECT * FROM people WHERE id=$1", id)

	person := new(Person)
	err := row.Scan(&person.Id, &person.FirstName, &person.Secondname, &person.Mail, &person.Gender)
	if err != nil {
		return nil, err
	}
	return person, err

}
