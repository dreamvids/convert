package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
		file, err := os.Create(TempDir + strconv.Itoa(id) + ".video")
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		io.Copy(file, r.Body)
		file.Close()

		c1 := NewConversion(id, FormatWebM, Resolution360p, StatusError)
		c2 := NewConversion(id, FormatWebM, Resolution720p, StatusError)

		c3 := NewConversion(id, FormatMp4, Resolution360p, StatusError)
		c4 := NewConversion(id, FormatMp4, Resolution720p, StatusError)

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

		err = DatabaseInsertConversion(&c3)
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		err = DatabaseInsertConversion(&c4)
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		cs = make([]Conversion, 4)
		cs[0] = c1
		cs[1] = c2
		cs[2] = c3
		cs[3] = c4

		err = cs[0].Start()
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		err = cs[1].Start()
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		err = cs[2].Start()
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		err = cs[3].Start()
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	err = json.NewEncoder(w).Encode(cs)
	if err != nil {
		SendErr(w, r, http.StatusInternalServerError, err)
		return
	}
}

func SendErr(w http.ResponseWriter, r *http.Request, status int, err error) {
	log.Println("Error", r.RemoteAddr, r.Method, r.RequestURI, ":", err)

	msg := strings.Replace(err.Error(), "\"", "\\\"", -1)

	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"failed\": true, \"message\": \"%s\"}", msg)
}
