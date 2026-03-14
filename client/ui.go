package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

func startUI(conn net.Conn, serverScanner *bufio.Scanner, message string) {
	app := tview.NewApplication()

	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	messages := tview.NewTextView().
	    SetScrollable(true).
	    SetDynamicColors(true).
	    SetChangedFunc(func() { app.Draw() })

	input := tview.NewInputField().
	    SetLabel("[white::b]>[-] ")

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(messages, 0, 1, false).
		AddItem(input, 1, 0, true)

	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		message := strings.TrimSpace(input.GetText())
		if message == "" {
			return
		}

		fmt.Fprintln(conn, message)
		input.SetText("")
	})

	go func() {
		update(app, messages, message)

		for serverScanner.Scan() {
			message := serverScanner.Text()

			// TODO: fix server disconnect message
			// ...and somehow try fixing terminal breaking?
			// (only happens when clients are on server after it disconnected)
			if message == "SERVER_DISCONNECT" {
				fmt.Fprintf(messages, "\nServer disconnected!")
				os.Exit(0)
			}
			update(app, messages, message)
		}
		os.Exit(0)
	}()

	if err := app.SetRoot(layout, true).Run(); err != nil {
		fmt.Println(err)
	}
}

func update(app *tview.Application, messages *tview.TextView, message string) {
	app.QueueUpdateDraw(func() {
		fmt.Fprintf(messages, "%s\n", message)
		messages.ScrollToEnd()
	})
}
