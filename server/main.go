package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
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
	history []string

	mut sync.RWMutex
	isShutdown bool
}

// TODO: add message encryption

func main() {
	f := parseFlags()

	listener, err := net.Listen("tcp", ":" + f.port)
	if err != nil {
		errorMessage(fmt.Sprintf("Port %s is already in use", f.port))
		os.Exit(1)
	}

	if (!f.verbose || f.logPath != "") && os.Getenv("SERVER_CHILD") != "1" {
		fmt.Printf("\nServer started on localhost:%s\n–> run `make kill` to stop it\n", f.port)

	    cmd := exec.Command(os.Args[0], os.Args[1:]...)
	    cmd.Env = append(os.Environ(), "SERVER_CHILD=1")

	    cmd.Stdout = nil
	    cmd.Stderr = nil
	    cmd.Stdin = nil

	    if err := cmd.Start(); err != nil {
			errorMessage(fmt.Sprintf("Unable to start server: %s", err))
	        os.Exit(1)
	    }
	    os.Exit(0)
	}

	setupLogging(f)
	defer listener.Close()

	fmt.Println()
	log.Printf("Server started on localhost:%s", f.port)

	server := &Server {
		clients: make(map[net.Conn]string),
		broadcast: make(chan Message),
	}

	go server.handleBroadcast()
	go server.handleConnections(listener)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit

	if f.verbose {
  		fmt.Print("\r\033[K")
	}
	server.shutdown(f.port)
}

func (s *Server) handleConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()

		if err != nil && !s.isShutdown {
			log.Printf("Error connecting to the server: %s", err)
			return
		}

		if conn == nil {
			return
		}

		log.Printf("New connection: %s", conn.RemoteAddr())
		go s.handleClient(conn)
	}
}

func (s *Server) handleBroadcast() {
	for message := range s.broadcast {
	 	if message.sender != nil {
            s.mut.Lock()
            s.history = append(s.history, message.content)

            if len(s.history) > MAX_HISTORY {
                s.history = s.history[1:]
            }

            s.mut.Unlock()
        }

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

func (s *Server) shutdown(port string) {
	s.isShutdown = true

	for conn := range s.clients {
		fmt.Fprintln(conn, "SERVER_DISCONNECT")
		conn.Close()
	}

	log.Printf("Shutting down server on localhost:%s", port)
	os.Exit(0)
}
