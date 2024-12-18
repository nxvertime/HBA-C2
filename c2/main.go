package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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
type ReqHeartBeat struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// //////// COMMANDS
type Command struct {
	Name        string
	Args        []string
	NoMaxArgs   bool
	Description string
}

// //////// LISTS ///////////////
// //// AVAILABLE COMMANDS
var commands = []Command{}

// //////////// WEB SERVER CONF ////////////
const PORT = ":443"

// ///////////  LOGS / CONSOLE MANAGEMENT ////////////////
var l = log.New(os.Stdout, "", 0)
var reader = bufio.NewReader(os.Stdin)
var customPrefix = "[LOG]"
var inputPrefix = "# "

func NewCommand(name string, args []string, noMaxArgs bool, description string) Command {
	if !noMaxArgs {
		noMaxArgs = true
	}
	return Command{name, args, noMaxArgs, description}
}
func InitCommands() {
	commands = append(commands, NewCommand("help", []string{}, false, "Show help"))
	commands = append(commands, NewCommand("version", []string{}, false, "Show version"))
	commands = append(commands, NewCommand("show", []string{"sessions / sessiond_id"}, false, "Show session list"))
	commands = append(commands, NewCommand("exec", []string{"session_id", "Args"}, true, "Execute a command on a specified session"))
	commands = append(commands, NewCommand("shell", []string{"session_id"}, false, "Start a remote shell session on a specified session"))

}

func Help() {
	fmt.Println("---- HELP ----")
	for i := 0; i < len(commands); i++ {
		fmt.Println("Usage:")
		fmt.Print("  " + commands[i].Name)
		for a := 0; a < len(commands[i].Args); a++ {
			fmt.Print(" ")
			fmt.Print(commands[i].Args[a])
		}
		fmt.Print("\n")
		fmt.Println("Description: ")
		fmt.Println("  " + commands[i].Description)

	}

}

func Log(l *log.Logger, msg string) {
	l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + customPrefix + " "))
	l.Print(msg)
	fmt.Print(inputPrefix)
}

func Error(l *log.Logger, msg string) {
	l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + "[ERROR]" + " "))
	l.Print(msg)

}

func UserInput(r *bufio.Reader) {
	for {
		fmt.Print(inputPrefix)
		input, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)

		Interpreter(l, input)
	}
}

func Interpreter(l *log.Logger, input string) {
	if input == "" {
		return
	}
	splittedInput := strings.Split(input, " ")

	cmdName := splittedInput[0]

	command := NewCommand("", []string{""}, false, "")
	for i := 0; i < len(commands); i++ {

		if (commands[i].Name) == cmdName {

			command = commands[i]
		}

	}
	if command.Name == "" {
		Error(l, "Command not found, type help")
		return
	}

	if !(len(command.Args) == (len(splittedInput) - 1)) {

		Error(l, "Invalid number of arguments given, "+strconv.Itoa(len(command.Args))+" expected")
		Log(l, "Usage: ")
		Log(l, "  "+command.Name+" "+strings.Join(command.Args, " "))

		fmt.Print("\n")
		return
	}

	switch command.Name {
	case "help":
		Help()

	}

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

func HeartBeat(w http.ResponseWriter, req *http.Request) {

}

func main() {
	InitCommands()
	go UserInput(reader)
	http.HandleFunc("/helloworld", HelloWorld)
	http.HandleFunc("/getSID", GetSID)
	Log(l, "Starting the webserver on "+PORT)

	err := http.ListenAndServeTLS(PORT, "certs/server.crt", "certs/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	defer Log(l, "Stopping the webserver...")

}
