package main

import (
	"bufio"
	"github.com/rivo/tview"
	"log"
	"os"
	"time"
)

var l = log.New(os.Stdout, "", 0)
var reader = bufio.NewReader(os.Stdin)
var customPrefix = "[DEBUG]"
var inputPrefix = "==> "

func LogEx(msg string, displayInputPrefix bool) {
	" " + "[blue::b]LOG[-::-]"
}

func Log(msg string) {
	LogEx(msg, false)
}
func DbgMsgEx(msg string, displayInputPrefix bool) {
	if !*verbose {
		return
	}
	UILog(textView, time.Now().Format("2006-01-02 15:04:05")+(" "+"[darkcyan::b]DBG[-::-]"+" ")+(msg))
}

func DbgMsg(msg string) {
	DbgMsgEx(msg, false)
}

func Error(msg string) {
	UILog(textView, time.Now().Format("2006-01-02 15:04:05")+(" "+"[red::b]ERROR[-::-]"+" ")+tview.Escape(msg))
}
