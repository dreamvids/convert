package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func HandleConvert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendErr(w, r, http.StatusBadRequest, fmt.Errorf("Invalid video id"))
		return
	}

	e, err := DatabaseVideoExists(id)
	if err != nil {
		SendErr(w, r, http.StatusInternalServerError, err)
		return
	}

	if !e {
		SendErr(w, r, http.StatusNotFound, fmt.Errorf("Video not found"))
		return
	}

	cs, err := DatabaseGetVideoConversions(id)
	if err != nil {
		SendErr(w, r, http.StatusInternalServerError, err)
		return
	}

	if len(cs) == 0 {
		c1 := NewConversion(id, FormatWebM, StatusConverting)
		c2 := NewConversion(id, FormatMp4, StatusConverting)

		err = DatabaseInsertConversion(&c1)
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		err = DatabaseInsertConversion(&c2)
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func SendErr(w http.ResponseWriter, r *http.Request, status int, err error) {
	log.Println("Error", r.RemoteAddr, r.Method, r.RequestURI, ":", err)

	msg := strings.Replace(err.Error(), "\"", "\\\"", -1)

	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"failed\": true, \"message\": \"%s\"}", msg)
}
