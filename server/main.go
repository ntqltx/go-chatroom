package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	clients map[net.Conn]string
	broadcast chan string
	mut sync.Mutex
}

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server started on :8080")

	server := &Server{
		clients: make(map[net.Conn]string),
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting to the server:", err)
			return
		}

		// fmt.Println("New connection:", conn)
		go server.handleClient(conn)
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

	fmt.Printf("%s joined the chat!\n", username)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("%s: %s\n", username, message)
	}

	s.mut.Lock()
	delete(s.clients, conn)
	s.mut.Unlock()

	fmt.Printf("%s left the chat!\n", username)

	// fmt.Println("Client disconnected:", conn.RemoteAddr())
}
