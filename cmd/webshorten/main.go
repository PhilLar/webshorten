package main

import (
	"encoding/json"
	"flag"
	//"fmt"
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

var flagPort *int = flag.Int("port", 0, "use '-port=NUMBER' flag to run server on specific port(default: 5000)")
const defaultPort string = ":5000"

// ListenAndServeWrapper is a wrapper for http.ListenAndServ to let one specidfy PORT
func ListenAndServeWrapper(a func(addr string, handler http.Handler) error) error {
	flag.Parse()
	if *flagPort > 0 {
		return a(":"+strconv.Itoa(*flagPort), nil)
	} else if path, exists := os.LookupEnv("PORT"); exists {
		return a(":"+path, nil)
	}
	return a(defaultPort, nil)
}

func main() {
	http.HandleFunc("/api/v1/shorten", shortenURL)
	err := ListenAndServeWrapper(http.ListenAndServe)
	if err != nil {
		log.Fatal(err)
	}
}
