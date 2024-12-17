package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

//////////// STRUCTS DECLARATION ////////////

// /////// RESPONSES
type ResGetSID struct {
	SessionId  string `json:"sessionId"`
	WelcomeMsg string `json:"welcomeMsg"`
}

// //////////// WEB SERVER CONF ////////////
const PORT = ":443"

// ///////////  LOGS / CONSOLE MANAGEMENT ////////////////
var l = log.New(os.Stdout, "", 0)
var customPrefix = "[LOG]"

func Log(l *log.Logger, msg string) {
	l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + (" " + customPrefix + " "))
	l.Print(msg)
}

// /////// OTHER THINGS ///////////
const CHARS = "abcdefghijklmnopqrstuvwxyz123456789"

func createSID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHARS[rand.Intn(len(CHARS))]

	}
	return string(b)
}

// /////// ROUTES PROCESSING /////////
func HelloWorld(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello World!\n"))
}

func GetSID(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	ipHeader := req.RemoteAddr
	sid := createSID(10)
	Log(l, ("From " + ipHeader + "=> GET /getSID SID: " + sid))
	res := ResGetSID{SessionId: sid, WelcomeMsg: "Welcome aboard:)"}

	jsonData, err := json.Marshal(res)
	if err != nil {
		Log(l, ("Error JSON marshalling: " + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

	}
	w.Write([]byte(jsonData))
}

func main() {

	http.HandleFunc("/helloworld", HelloWorld)
	http.HandleFunc("/getSID", GetSID)
	Log(l, "Starting the webserver on "+PORT)

	err := http.ListenAndServeTLS(PORT, "certs/server.crt", "certs/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	defer Log(l, "Stopping the webserver...")

}
