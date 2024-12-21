package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

var l = log.New(os.Stdout, "", 0)
var reader = bufio.NewReader(os.Stdin)
var customPrefix = "[LOG]"
var inputPrefix = "# "

func LogEx(l *log.Logger, msg string, displayInputPrefix bool) {
	l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + customPrefix + " "))
	l.Print(msg)
	if displayInputPrefix {
		fmt.Print(inputPrefix)

	}
}

func Log(l *log.Logger, msg string) {
	LogEx(l, msg, false)
}

func Error(l *log.Logger, msg string) {
	l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + "[ERROR]" + " "))
	l.Print(msg)

}
