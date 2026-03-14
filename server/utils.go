package main

import (
	"net"
	"sync"
	"github.com/fatih/color"
)

type Message struct {
	sender net.Conn
	content string
}

type Server struct {
	clients map[net.Conn]string
	broadcast chan Message
	mut sync.RWMutex
}

// -- slice with all possible username colors
var colorList = []*color.Color{
	color.New(color.FgRed),
	color.New(color.FgGreen),
	color.New(color.FgYellow),
	color.New(color.FgBlue),
	color.New(color.FgMagenta),
	color.New(color.FgCyan),
}

// -- timestamp coloring
var timestampStyle = color.New(color.FgWhite, color.Faint, color.Bold)

// -- function to get random username color
func getUserColor(username string) *color.Color {
	total := 0
	for _, c := range username {
		total += int(c)
	}
	return colorList[total % len(colorList)]
}
