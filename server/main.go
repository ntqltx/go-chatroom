package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server started on :8080")

	server := &Server {
		clients: make(map[net.Conn]string),
		broadcast: make(chan Message),
	}
	go server.handleBroadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting to the server:", err)
			return
		}

		fmt.Println("New connection:", conn)
		go server.handleClient(conn)
	}
}

func (s *Server) handleBroadcast() {
	for message := range s.broadcast {
        s.mut.RLock()
        for conn := range s.clients {
        	if conn != message.sender {
            	fmt.Fprintln(conn, message.content)
         	}
        }
        s.mut.RUnlock()
    }
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	scanner.Scan()
	username := scanner.Text()

	s.mut.Lock()
	s.clients[conn] = username
	s.mut.Unlock()

	var nameColor = getUserColor(username)
	colorUsername := nameColor.Sprint(username)
	fmt.Fprintf(conn, "Connected as %s!\n", colorUsername)

	s.broadcast <- Message {
		sender: conn,
		content: fmt.Sprintf("%s joined!", colorUsername),
	}

	for scanner.Scan() {
		message := scanner.Text()
		timestamp := timestampStyle.Sprintf("[%s]:", time.Now().Format("15:04"))

		var formatted = fmt.Sprintf(
			"%s %s: %s", timestamp,
			colorUsername, message,
		)
		fmt.Fprintln(conn, formatted)

		s.broadcast <- Message {
			sender: conn,
			content: formatted,
		}
	}

	s.mut.Lock()
	delete(s.clients, conn)
	s.mut.Unlock()

	s.broadcast <- Message {
		sender: conn,
		content: fmt.Sprintf("%s left", colorUsername),
	}
	fmt.Println("Client disconnected:", conn.RemoteAddr())
}
