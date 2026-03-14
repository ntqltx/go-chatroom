package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// -- connecting to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// -- username getter
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Username: ")
	scanner.Scan()

	username := scanner.Text()
	fmt.Fprintln(conn, username)

	// -- receive messages
	go func() {
		serverScanner := bufio.NewScanner(conn)
		for serverScanner.Scan() {
			fmt.Printf("\r\033[K%s\n> ", serverScanner.Text())
		}

		fmt.Print("\n\033[A\033[K")
		fmt.Println("Server disconnected!")
		os.Exit(0)
	}()

	// -- message loop
 	for {
        scanner.Scan()
        message := scanner.Text()

        if message == "" {
        	continue
        }

        fmt.Print("\033[A\033[K")
        fmt.Fprintln(conn, message)
    }
}
