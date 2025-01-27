package main

import "time"

//////////// STRUCTS DECLARATION ////////////

// /////// RESPONSES
type ResGetSID struct {
	SessionId  string `json:"sessionId"`
	WelcomeMsg string `json:"welcomeMsg"`
}

type ResHeartBeat struct {
	SessionId string        `json:"sessionId"`
	Type      string        `json:"type"`
	Args      []interface{} `json:"args"`
}

type ResRegister struct {
	ResMsg string `json:"resMsg"`
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
type InputCommand struct {
	Name        string
	Args        []string
	NoMaxArgs   bool
	Description string
}

type Command struct {
	Type string
	Args []string
}

// /////////// ZOMBIES
type Zombie struct {
	Id            int
	SessionId     string
	RemoteAddr    string
	RemotePort    string
	UserName      string
	Country       string
	FirstConnTime time.Time
	LastConnTime  time.Time
}
