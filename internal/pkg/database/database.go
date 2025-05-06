package database

import (
	"database/sql"
	"fmt"
)

type PostgresDB struct {
	DB *sql.DB
}

func InitDB(user, password, host, port, dbname string) PostgresDB {
	conn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic("database does not work")
	}
	return PostgresDB{DB: db}
}
