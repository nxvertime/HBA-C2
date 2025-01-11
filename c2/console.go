package main

import (
	"bufio"
	"log"
	"os"
	"time"
)

var l = log.New(os.Stdout, "", 0)
var reader = bufio.NewReader(os.Stdin)
var customPrefix = "[DEBUG]"
var inputPrefix = "==> "

func LogEx(msg string, displayInputPrefix bool) {

	//l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + "[LOG]" + " "))
	//l.Print(msg)
	//if displayInputPrefix {
	//	fmt.Print(inputPrefix)

	//}

	UILog(textView, time.Now().Format("2006-01-02 15:04:05")+(" "+"[LOG]"+" ")+msg)
}

func Log(msg string) {
	LogEx(msg, false)
}
func DbgMsgEx(msg string, displayInputPrefix bool) {
	if !*verbose {
		return
	}
	//l.SetPrefix("\033[2K\r" + time.Now().Format("2006-01-02 15:04:05") + (" " + "[DBG]" + " "))
	//l.Print(msg)
	//if displayInputPrefix {
	//	fmt.Print(inputPrefix)

	//}
	UILog(textView, time.Now().Format("2006-01-02 15:04:05")+(" "+"[DBG]"+" ")+msg)

}

func DbgMsg(msg string) {
	DbgMsgEx(msg, false)
}

func Error(msg string) {
	//l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + (" " + "[ERROR]" + " "))
	//l.Print(msg)
	UILog(textView, time.Now().Format("2006-01-02 15:04:05")+(" "+"[ERROR]"+" ")+msg)

}
