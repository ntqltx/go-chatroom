package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// TODO: add address inputing
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// TODO:
	// 1. clean username input after entering
	// 2. check if username is free
	// 3. check if username isn't more than 20 symbols
	//
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Username: ")
	scanner.Scan()

	username := scanner.Text()
	fmt.Fprintln(conn, username)

	go func() {
		serverScanner := bufio.NewScanner(conn)
		for serverScanner.Scan() {
			fmt.Printf("\r\033[K%s\n> ", serverScanner.Text())
		}

		// TODO: make server disconnected message on server-side
		fmt.Print("\n\033[A\033[K")
		fmt.Println("Server disconnected!")
		os.Exit(0)
	}()

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
