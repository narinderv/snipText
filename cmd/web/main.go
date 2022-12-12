package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Common for functions across the package
type configuration struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// Command Line Arguments
	serverAddr := flag.String("addr", ":8888", "Network address")
	flag.Parse()

	config := &configuration{
		infoLog:  log.New(os.Stdout, "INFO|", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR|", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// HTTP Server
	httpServer := &http.Server{
		Addr:     *serverAddr,
		Handler:  config.routes(),
		ErrorLog: config.errorLog,
	}

	config.infoLog.Printf("Starting the server on port %s", *serverAddr)
	err := httpServer.ListenAndServe()
	config.errorLog.Fatal(err)
}
