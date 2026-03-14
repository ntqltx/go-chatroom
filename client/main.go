package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// -- connecting to the server
	// TODO: disconnect user if server closes
	//
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
	fmt.Printf("Connected as %s!\n", username)

	// -- message loop
 	for {
        fmt.Print("> ")
        scanner.Scan()
        message := scanner.Text()
        fmt.Fprintln(conn, message)
    }
}
