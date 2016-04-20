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
		path := TempDir + strconv.Itoa(id) + ".video"
		file, err := os.Create(path)
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		io.Copy(file, r.Body)
		file.Close()

		info, err := ProbeVideo(path)
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}

		css := make([]*Conversion, 2)
		css[0] = NewConversion(id, FormatWebM, Resolution360p, StatusError)
		css[1] = NewConversion(id, FormatMp4, Resolution360p, StatusError)

		if info.Width >= 1280 {
			css = append(css, NewConversion(id, FormatWebM, Resolution720p, StatusError))
			css = append(css, NewConversion(id, FormatMp4, Resolution720p, StatusError))
		}

		for _, c := range css {
			err = DatabaseInsertConversion(c)
			if err != nil {
				SendErr(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		for _, c := range css {
			err = c.Start()
			if err != nil {
				SendErr(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		err = json.NewEncoder(w).Encode(css)
		if err != nil {
			SendErr(w, r, http.StatusInternalServerError, err)
			return
		}
	} else {
		err = json.NewEncoder(w).Encode(cs)
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
