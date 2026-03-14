package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
)

const CLEAN_LINE = "\033[A\033[K"

func main() {
	// TODO: add address inputing
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	defer conn.Close()
	scanner := bufio.NewScanner(os.Stdin)

	var username string
	for {
		fmt.Print("Enter Username: ")
		scanner.Scan()

		username = strings.TrimSpace(scanner.Text())
		fmt.Print(CLEAN_LINE)

		if err := validateUsername(username); err != "" {
			errorMessage(err)
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

		// temp
		fmt.Printf("\033[2J\033[H")
		fmt.Printf("\r\033[K%s\n> ", response)

		go func() {
			for serverScanner.Scan() {
				message := serverScanner.Text()
				fmt.Printf("\r\033[K%s\n> ", message)
			}

			// TODO: make server disconnected message on server-side
			fmt.Printf("\n%s", CLEAN_LINE)
			os.Exit(0)
		}()

		break
	}

	// TODO: fix the prompt line
 	for {
        scanner.Scan()
        message := strings.TrimSpace(scanner.Text())

        if message == "" {
       		fmt.Printf("%s> ", CLEAN_LINE)
        	continue
        }

        fmt.Print(CLEAN_LINE)
        fmt.Fprintln(conn, message)
    }
}

func validateUsername(username string) string {
	l := len(username)
	switch {
		case l == 0: return "Username cannot be empty"
		case l > 20: return "Username cannot be more than 20 characters"
	}
	return ""
}

func errorMessage(content string) {
    errorStyle := color.New(color.FgRed, color.Bold)
    fmt.Println(errorStyle.Sprint("[Error] ")+content)
}
