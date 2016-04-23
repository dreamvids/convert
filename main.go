package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	Version = "0.0.1"
)

var (
	TempDir = "temp/"
)

func main() {
	log.Println("Starting convert version", Version)

	tempDir := os.Getenv("DV_CONVERT_TEMP")
	if len(tempDir) > 0 {
		TempDir = tempDir
	}

	dbHost := os.Getenv("DV_DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "db"
	}

	dbPort, err := strconv.Atoi(os.Getenv("DV_DB_PORT"))
	if err != nil || dbPort == 0 {
		dbPort = 3306
	}

	dbUser := os.Getenv("DV_DB_USER")
	if len(dbUser) == 0 {
		dbUser = "root"
	}

	dbPass := os.Getenv("DV_DB_PASSWORD")
	if len(dbPass) == 0 {
		dbPass = "root"
	}

	dbName := os.Getenv("DV_DB_NAME")
	if len(dbName) == 0 {
		dbName = "dreamvids"
	}

	err = os.MkdirAll(TempDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = DatabaseInit(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal("Database initialization: ", err)
	}

	defer Database.Close()

	r := mux.NewRouter()

	r.HandleFunc("/convert/{id}", HandleConvert)
	http.Handle("/", r)

	log.Println("Listening to 0.0.0.0:8001...")
	log.Fatal(http.ListenAndServe(":8001", nil))
}
