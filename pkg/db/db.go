package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func InitDB(host, port, user, password, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %w", err)
	}

	return db, nil
}
