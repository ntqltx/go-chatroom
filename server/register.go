package main

import (
	"net"
	"fmt"
	"bufio"
)

func (s *Server) registerClient(conn net.Conn) (string, string) {
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return "", ""
	}

	username := scanner.Text()
	colorUsername := getUserColor(username).Sprint(username)

	s.systemBroadcast(fmt.Sprintf("%s joined!", colorUsername), conn)

	s.mut.Lock()
	s.clients[conn] = username
	s.mut.Unlock()

	fmt.Fprintf(conn, "Connected as %s\n", colorUsername)
	fmt.Printf("%s connected\n", username)

	return username, colorUsername
}

func (s *Server) unregisterClient(conn net.Conn, username, colorUsername string) {
	s.mut.Lock()
	delete(s.clients, conn)
	s.mut.Unlock()

	s.systemBroadcast(fmt.Sprintf("%s left!", colorUsername), conn)
	fmt.Printf("%s disconnected\n", username)
}
