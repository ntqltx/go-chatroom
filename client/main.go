package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// TODO: add address input
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	serverScanner, message, err := login(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	startUI(conn, serverScanner, message)
}

func login(conn net.Conn) (*bufio.Scanner, string, error) {
	scanner := bufio.NewScanner(os.Stdin)

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

		serverScanner := bufio.NewScanner(conn)
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
