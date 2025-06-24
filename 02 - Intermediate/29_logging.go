package main

import (
	"log"
	"os"
)

func main() {
	log.Println("This is a log message")

	log.SetPrefix("INFO: ")
	log.Println("This is an INFO flag")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("this is a log message with date, time")

	infoLogger.Println("Using info logger")
	warnLogger.Println("Using warn logger")
	errorLogger.Println("Using error logger")

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		errorLogger.Fatalf("Failed to open log file: %s", err)
	}
	defer file.Close()

	debugLogger := log.New(file, "DEBUG: ", log.Ldate|log.Ltime)
	debugLogger.Println("Created a logger file")
}

var (
	infoLogger  = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)
	warnLogger  = log.New(os.Stdout, "WARN: ", log.Ldate|log.Lshortfile)
	errorLogger = log.New(os.Stdout, "Error: ", log.Ldate|log.Lshortfile)
)
