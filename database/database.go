package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect(dbConn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
