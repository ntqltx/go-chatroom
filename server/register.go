package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func (s *Server) registerClient(conn net.Conn) (string, string) {
	scanner := bufio.NewScanner(conn)

	for {
		if !scanner.Scan() {
			return "", ""
		}
		username := strings.TrimSpace(scanner.Text())

		s.mut.RLock()
		taken := false

		for _, name := range s.clients {
			if name == username {
				taken = true
				break
			}
		}

		s.mut.RUnlock()

		if taken {
			fmt.Fprintln(conn, "USERNAME_TAKEN")
			continue
		}

		colored := colorUsername(username)
		s.systemBroadcast(fmt.Sprintf("%s joined!", colored), conn)

		s.mut.Lock()
		s.clients[conn] = username
		s.mut.Unlock()

		s.mut.RLock()
		for _, msg := range s.history {
		    fmt.Fprintln(conn, msg)
		}
		s.mut.RUnlock()

		fmt.Fprintf(conn, "Connected as %s\n", colored)
		fmt.Printf("%s connected\n", username)

		return username, colored
	}
}

func (s *Server) unregisterClient(conn net.Conn, username, colorUsername string) {
	s.mut.Lock()
	delete(s.clients, conn)
	s.mut.Unlock()

	s.systemBroadcast(fmt.Sprintf("%s left!", colorUsername), conn)
	fmt.Printf("%s disconnected\n", username)
}
