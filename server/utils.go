package main

import (
	"fmt"
	"net"
	"time"

	"github.com/fatih/color"
)

const CLEAN_LINE = "\033[A\033[K"
var errorStyle = color.New(color.FgRed, color.Bold)

// -- slice with all possible username color tags
var colorList = []string{
	"[red]",
	"[green]",
	"[yellow]",
	"[blue]",
	"[purple]",
	"[cyan]",
}

func getUserColor(username string) string {
	total := 0
	for _, c := range username {
		total += int(c)
	}
	return colorList[total % len(colorList)]
}

func colorUsername(username string) string {
	return fmt.Sprintf("%s%s[-]", getUserColor(username), username)
}

func (s *Server) formatMessage(username, message string) string {
	timestamp := fmt.Sprintf("[gray]%s[-]", time.Now().Format("[15:04]"))
	return fmt.Sprintf("%s %s: %s", timestamp, colorUsername(username), message)
}

func (s *Server) systemBroadcast(content string, target net.Conn) {
	s.broadcast <- Message{sender: nil, target: target, content: content}
}

func errorMessage(content string) {
	fmt.Print(CLEAN_LINE)
	fmt.Println(errorStyle.Sprint("[Error] ") + content)
}
