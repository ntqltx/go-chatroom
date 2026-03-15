package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type flags struct {
	verbose bool
	port string
	logPath string
}

func parseFlags() *flags {
	f := &flags{}

	flag.BoolVar(&f.verbose, "verbose", false, "show server logs in terminal")
	flag.StringVar(&f.port, "port", "8080", "port to listen on")
	flag.StringVar(&f.logPath, "log", "", "save logs to file")

	flag.Parse()

	if _, err := strconv.Atoi(f.port); err != nil {
		errorMessage("Port must be a number")
		os.Exit(1)
	}

	return f
}

func setupLogging(f *flags) {
	writers := []io.Writer{io.Discard}

	if f.logPath != "" {
		file, err := os.OpenFile(f.logPath, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0644)

		if err != nil {
			errorMessage(fmt.Sprintf("Cannot open log file: %s", err))
			os.Exit(1)
		}

		writers = []io.Writer{file}
	}

	if f.verbose {
		writers = append(writers, os.Stdout)
	}

	log.SetOutput(io.MultiWriter(writers...))
	log.SetFlags(log.Ltime)
}
