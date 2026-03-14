package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Message struct {
	sender net.Conn
	target net.Conn
	content string
}

type Server struct {
	clients map[net.Conn]string
	broadcast chan Message
	mut sync.RWMutex
	isShutdown bool
}

func main() {
	// TODO: add address input
	// ALSO ADD ENCRYPTION
	//
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server started on localhost:8080")

	server := &Server {
		clients: make(map[net.Conn]string),
		broadcast: make(chan Message),
	}

	go server.handleBroadcast()
	go server.handleConnections(listener)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit

	fmt.Print("\r\033[K")
	server.shutdown()
}

func (s *Server) handleConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()

		if err != nil && !s.isShutdown {
			fmt.Println("Error connecting to the server:", err)
			return
		}

		if conn == nil {
			return
		}

		fmt.Println("New connection:", conn.RemoteAddr())
		go s.handleClient(conn)
	}
}

func (s *Server) handleBroadcast() {
	for message := range s.broadcast {
		s.mut.RLock()

		for conn := range s.clients {
			if conn == message.sender || conn == message.target {
                continue
            }

			switch message.sender {
				case nil:
					formatted := fmt.Sprintf("[white::b][Server[][-:-:-] %s", message.content)
					fmt.Fprintln(conn, formatted)

				default: fmt.Fprintln(conn, message.content)
			}
		}

		s.mut.RUnlock()
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	username, colorUsername := s.registerClient(conn)
	if username == "" {
		return
	}

	defer s.unregisterClient(conn, username, colorUsername)
	s.receiveMessages(conn, colorUsername)
}

func (s *Server) receiveMessages(conn net.Conn, colorUsername string) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		formatted := s.formatMessage(colorUsername, message)

		fmt.Fprintln(conn, formatted)
		s.broadcast <- Message{sender: conn, content: formatted}
	}
}

func (s *Server) shutdown() {
	s.isShutdown = true

	for conn := range s.clients {
		fmt.Fprintln(conn, "SERVER_DISCONNECT")
		conn.Close()
	}

	fmt.Println("\nShutting down...")
}
