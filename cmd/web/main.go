package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/narinderv/snipText/pkg/models"
	"github.com/narinderv/snipText/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

// Common for functions across the package
type configuration struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	sessionManager *sessions.Session
	templateCache  map[string]*template.Template
	snips          interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.SnipText, error)
		GetLatest() ([]*models.SnipText, error)
	}
	users interface {
		Insert(string, string, string) error
		Get(id int) (*models.User, error)
		Authenticate(string, string) (int, error)
	}
}

func main() {

	// Command Line Arguments
	serverAddr := flag.String("addr", ":8888", "Network address")
	dbDetails := flag.String("conn", "web:sniptext@/sniptext?parseTime=true",
		"Database connnection detail (user:password@/database-name)?parseTime=true")
	sessionKey := flag.String("key", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Session Secret Key (32 bit)")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO|", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR|", log.Ldate|log.Ltime|log.Lshortfile)

	dbConnection, err := connectToDatabase(*dbDetails)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer dbConnection.Close()

	dbConnUser, err := connectToDatabase(*dbDetails)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer dbConnUser.Close()

	tmplCache, err := templateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Session Manager
	sessionManager := sessions.New([]byte(*sessionKey))
	sessionManager.Lifetime = time.Hour * 12

	config := &configuration{
		infoLog:        infoLog,
		errorLog:       errorLog,
		sessionManager: sessionManager,
		templateCache:  tmplCache,
		snips:          &mysql.SnipModel{DB: dbConnection},
		users:          &mysql.UserModel{DB: dbConnUser},
	}

	// TLS specific configurations
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// HTTP Server
	httpServer := &http.Server{
		Addr:         *serverAddr,
		Handler:      config.routes(),
		ErrorLog:     config.errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting the server on port %s", *serverAddr)
	err = httpServer.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
