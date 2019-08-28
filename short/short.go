package short

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	//"fmt"
	"io/ioutil"
	//"log"
	"net/http"
	"net/url"
	"sync"
)

const CleanUrlApi string = "https://cleanuri.com/api/v1/shorten"
const RelinkApi string = "https://rel.ink/api/links/"

var flagRelink *bool = flag.Bool("relink", false, "use rel.ink service to shorten URL")

type CleanUrlAnswer struct {
	ResultUrl string `json:"result_url"`
	Error     string
}

// Wrapper for CleanUrl() and Relink()
func RunInParallel(a func(urlLink string) (string, error), args []string) map[string]error {
	shorts := make(map[string]error)
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(args))
	for _, arg := range args {
		go func(arg string) {
			defer wg.Done()
			short, err := a(arg)
			mutex.Lock()
			shorts[short] = err
			mutex.Unlock()
		}(arg)
	}
	wg.Wait()
	return shorts
}

func ShortenUrls() map[string]error {
	flag.Parse()
	if !*flagRelink {
		return RunInParallel(CleanUrl, flag.Args())
	} else {
		return RunInParallel(Relink, flag.Args())
	}
}

func CleanUrl(urlLink string) (string, error) {
	resp, err := http.PostForm(CleanUrlApi, url.Values{
		"url": {urlLink},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	answer := &CleanUrlAnswer{}
	err = json.Unmarshal(body, answer)
	if err != nil {
		return "", err
	}
	if answer.Error != "" {
		return "", errors.New(answer.Error)
	} else {
		return answer.ResultUrl, nil
	}
}

func Relink(urlLink string) (string, error) {
	jsn, err := json.Marshal(map[string]string{
		"url": urlLink,
	})
	if err != nil {
		return "", err
	}
	resp, err := http.Post(RelinkApi, "application/json", bytes.NewBuffer(jsn))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	answer := make(map[string]interface{})
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return "", err
	}
	shortValue, ok := answer["hashid"]
	if !ok {
		return "", errors.New("hashid not found")
	}
	short, ok := shortValue.(string)
	if !ok {
		return "", errors.New("failed in type assertion")
	}

	cleanUrl := "https://rel.ink/" + short
	return cleanUrl, err
}

// func main() {
// 	resultMap := ShortenUrls()
// 	for short, err := range resultMap {
// 		if err == nil {
// 			fmt.Println(short)
// 		} else {
// 			log.Fatal(err)
// 		}
// 	}
// }
