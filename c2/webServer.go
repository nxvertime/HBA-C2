package main

import (
	"database/sql"
	"log"
	"net/http"
)

func StartWebServer(db *sql.DB) {
	http.HandleFunc("/helloworld", HelloWorld)
	http.HandleFunc("/getSID", GetSID)
	http.HandleFunc("/register", Register(*db))
	http.HandleFunc("/heartBeat", HeartBeat(*db))
	DbgMsgEx("Starting the webserver on "+PORT, true)

	err := http.ListenAndServeTLS(PORT, "certs/server.crt", "certs/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	defer DbgMsg("Stopping the webserver...")
}
