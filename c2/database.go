package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const DB_HOST = "localhost"
const DB_PORT = 3306
const DB_NAME = "hba_db"
const DB_USR = "hba_user"
const DB_PWD = "hbaPass123"

func InitDb() (*sql.DB, error) {
	var db *sql.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", DB_USR, DB_PWD, DB_HOST, DB_PORT, DB_NAME)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to DB : %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error verifying DB : %v", err)
	}

	return db, nil

}
