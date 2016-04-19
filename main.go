package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	Version = "0.0.1"
	TempDir = "temp/"
)

func main() {
	log.Println("Starting convert version", Version)

	err := os.MkdirAll(TempDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = DatabaseInit("127.0.0.1", 3306, "root", "root", "dreamvids")
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
