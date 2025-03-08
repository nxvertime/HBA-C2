package main

import (
	"database/sql"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application
var textView *tview.TextView
var inputField *tview.InputField
var flex *tview.Flex

func UILog(txtv *tview.TextView, msg string) {
	UILogEx(txtv, msg, true)
	textView.ScrollToEnd()
}

func UILogEx(txtv *tview.TextView, msg string, new_line bool) {
	//escapedMsg := tview.Escape(msg)
	if new_line {
		fmt.Fprintf(txtv, "%s\n", msg)

	} else {
		fmt.Fprintf(txtv, "%s", msg)

	}
	textView.ScrollToEnd()
}

func StartUI(db *sql.DB) {
	app = tview.NewApplication()
	textView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetTitle("HBA_CONSOLE").SetTitleAlign(tview.AlignLeft)

	textView.SetBorder(true)

	inputField = tview.NewInputField().
		SetLabel(inputPrefix).
		SetFieldWidth(80).
		SetFieldBackgroundColor(tcell.ColorBlack)
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := inputField.GetText()
			UILog(textView, inputPrefix+text)

			inputField.SetText("")
			Interpreter(text, db)
		}

	}).SetBorder(true)

	flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, false).
		AddItem(inputField, 3, 1, true)

	if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(inputField).Run(); err != nil {
		panic(err)
	}
}
