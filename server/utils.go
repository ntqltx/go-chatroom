package main

import (
	"fmt"
	"net"
	"time"

	"github.com/fatih/color"
)

// -- slice with all possible username colors
var colorList = []*color.Color{
	color.New(color.FgRed),
	color.New(color.FgGreen),
	color.New(color.FgYellow),
	color.New(color.FgBlue),
	color.New(color.FgMagenta),
	color.New(color.FgCyan),
}

// -- custom coloring
var systemStyle = color.New(color.FgHiWhite, color.Bold)
var timestampStyle = color.New(color.FgWhite, color.Faint, color.Bold)

func getUserColor(username string) *color.Color {
	total := 0
	for _, c := range username {
		total += int(c)
	}
	return colorList[total % len(colorList)]
}

func (s *Server) formatMessage(colorUsername, message string) string {
	timestamp := timestampStyle.Sprintf("[%s]", time.Now().Format("15:04"))
	return fmt.Sprintf("%s %s: %s", timestamp, colorUsername, message)
}

func (s *Server) systemBroadcast(content string, target net.Conn) {
    s.broadcast <- Message{sender: nil, target: target, content: content}
}
