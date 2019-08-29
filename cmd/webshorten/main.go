package main

import (
	"encoding/json"
	"flag"
	"github.com/PhilLar/webshorten/short"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
	log.Println(r.Method, r.URL)
	log.Println(r.Host)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

var port int

func init() {
	defPort := 5000
	if portVar, ok := os.LookupEnv("PORT"); ok {
		if portValue, err := strconv.Atoi(portVar); err == nil {
			defPort = portValue
		}
	}
	flag.IntVar(&port, "port", defPort, "port to listen on")
}

func main() {
	flag.Parse()
	http.HandleFunc("/api/v1/shorten", shortenURL)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
