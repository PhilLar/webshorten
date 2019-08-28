package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"short"
)

type ShortUrlWrapper struct {
	ShortUrl string `json:"shortURL"`
}

type LongUrlWrapper struct {
	LongUrl string `json:"longURL"`
}

func shortenUrl(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	longUrlWrapper := &LongUrlWrapper{}
	err = json.Unmarshal(body, longUrlWrapper)
	if err != nil {
		log.Fatal(err)
	}
	resUrl, err := short.CleanUrl(longUrlWrapper.LongUrl)
	if err != nil {
		log.Fatal(err)
	}
	shortUrlWrapper := &ShortUrlWrapper{ShortUrl: resUrl}
	js, err := json.Marshal(shortUrlWrapper)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(js)
	log.Println("hello function handler was executed")
}

func main() {
	http.HandleFunc("/api/v1/shorten", shortenUrl)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
