package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var conn net.Conn

	for {
		fmt.Print("Port (default 8080): ")
		scanner.Scan()

		port := strings.TrimSpace(scanner.Text())
		if port == "" { port = "8080" }

		var err error
		address := fmt.Sprintf("localhost:%s", port)
		conn, err = net.Dial("tcp", address)

		if err != nil {
			errorMessage(fmt.Sprintf("Could not connect to %s, try again!", address))
			continue
		}

		signal.Ignore(syscall.SIGTERM)
		break
	}
	defer conn.Close()

	serverScanner, message, err := login(conn, scanner)
	if err != nil {
		fmt.Println(err)
		return
	}
	startUI(conn, serverScanner, message)
}

func login(conn net.Conn, scanner *bufio.Scanner) (*bufio.Scanner, string, error) {
	serverScanner := bufio.NewScanner(conn)

	for {
		fmt.Print("Enter username: ")
		scanner.Scan()

		username := strings.TrimSpace(scanner.Text())
		fmt.Print(CLEAN_LINE)

		if errMsg := validateUsername(username); errMsg != "" {
			errorMessage(errMsg)
			continue
		}

		fmt.Fprintln(conn, username)
		serverScanner.Scan()
		response := serverScanner.Text()

		if response == "USERNAME_TAKEN" {
			errorMessage("Username is already taken, try another!")
			continue
		}

		fmt.Printf("\033[2J\033[H")
		return serverScanner, response, nil
	}
}
