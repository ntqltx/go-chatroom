package main

import (
	"fmt"

	"github.com/fatih/color"
)

const CLEAN_LINE = "\033[A\033[K"
var errorStyle = color.New(color.FgRed, color.Bold)

func validateUsername(username string) string {
	l := len(username)
	switch {
		case l == 0: return "Username cannot be empty"
		case l > 20: return "Username cannot be more than 20 characters"
	}
	return ""
}

func errorMessage(content string) {
	fmt.Println(errorStyle.Sprint("[Error] ") + content)
}
