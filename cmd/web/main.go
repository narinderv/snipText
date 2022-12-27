package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/narinderv/snipText/pkg/models/mysql"
)

// Common for functions across the package
type configuration struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snips         *mysql.SnipModel
	templateCache map[string]*template.Template
}

func main() {

	// Command Line Arguments
	serverAddr := flag.String("addr", ":8888", "Network address")
	dbDetails := flag.String("conn", "web:sniptext@/sniptext?parseTime=true",
		"Database connnection detail (user:password@/database-name)?parseTime=true")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO|", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR|", log.Ldate|log.Ltime|log.Lshortfile)

	dbConnection, err := connectToDatabase(*dbDetails)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer dbConnection.Close()

	tmplCache, err := templateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	config := &configuration{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snips:         &mysql.SnipModel{DB: dbConnection},
		templateCache: tmplCache,
	}

	// HTTP Server
	httpServer := &http.Server{
		Addr:     *serverAddr,
		Handler:  config.routes(),
		ErrorLog: config.errorLog,
	}

	infoLog.Printf("Starting the server on port %s", *serverAddr)
	err = httpServer.ListenAndServe()
	errorLog.Fatal(err)
}

func connectToDatabase(dbDetails string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbDetails)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
