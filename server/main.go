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

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting to the server:", err)
			return
		}

		fmt.Println("New connection:", conn)
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Received:", message)
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
}
