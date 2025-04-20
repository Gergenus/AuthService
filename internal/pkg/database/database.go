package database

import (
	"database/sql"
	"fmt"
)

type PostgresDB struct {
	DB *sql.DB
}

func InitDB(user, password, host, port, dbname, sslmode string) PostgresDB {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic("database does not work")
	}
	return PostgresDB{DB: db}
}
