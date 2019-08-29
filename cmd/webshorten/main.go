package main

import (
	"encoding/json"
	"github.com/PhilLar/webshorten/short"
	"io/ioutil"
	"log"
	"net/http"
)

type shortURLWrapper struct {
	ShortURL string `json:"ShortURL"`
}

type longURLWrapper struct {
	LongURL string `json:"LongURL"`
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	longURL := &longURLWrapper{}
	err = json.Unmarshal(body, longURL)
	if err != nil {
		log.Print("sdf", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	resURL, err := short.CleanURL(longURL.LongURL)
	if err != nil {
		log.Print("sdf", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	shortURL := &shortURLWrapper{ShortURL: resURL}
	js, err := json.Marshal(shortURL)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/api/v1/shorten", shortenURL)
	err := http.ListenAndServe(":5001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
