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
var customPrefix = "[DEBUG]"
var inputPrefix = "# "

func LogEx(l *log.Logger, msg string, displayInputPrefix bool) {

	l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + "[LOG]" + " "))
	l.Print(msg)
	if displayInputPrefix {
		fmt.Print(inputPrefix)

	}
}

func Log(l *log.Logger, msg string) {
	LogEx(l, msg, false)
}
func DbgMsgEx(l *log.Logger, msg string, displayInputPrefix bool) {
	if !*verbose {
		return
	}
	l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + "[DBG]" + " "))
	l.Print(msg)
	if displayInputPrefix {
		fmt.Print(inputPrefix)

	}
}

func DbgMsg(l *log.Logger, msg string) {
	DbgMsgEx(l, msg, false)
}

func Error(l *log.Logger, msg string) {
	l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + "[ERROR]" + " "))
	l.Print(msg)

}
