package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func New() *sql.DB {
	db, err := newDatabase()
	if err != nil {
		panic(err)
	}
	return db
}

func newDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func getConnectionString() string {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	connectingStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	return connectingStr
}
