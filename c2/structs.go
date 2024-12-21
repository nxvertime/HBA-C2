package main

import "time"

//////////// STRUCTS DECLARATION ////////////

// /////// RESPONSES
type ResGetSID struct {
	SessionId  string `json:"sessionId"`
	WelcomeMsg string `json:"welcomeMsg"`
}

type ResHeartBeat struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

// /////////// REQUESTS
type ReqRegister struct {
	SessionId string `json:"sessionId"`
}

type ReqHeartBeat struct {
	SessionId string `json:"sessionId"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// //////// COMMANDS
type Command struct {
	Name        string
	Args        []string
	NoMaxArgs   bool
	Description string
}

// /////////// ZOMBIES
type Zombie struct {
	SessionId     string
	RemoteAddr    string
	RemotePort    string
	UserName      string
	Country       string
	FirstConnTime time.Time
	LastConnTime  time.Time
}
