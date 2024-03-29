package short

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const cleanURLAPI string = "https://cleanuri.com/api/v1/shorten"
const relinkAPI string = "https://rel.ink/api/links/"

type cleanURLAnswer struct {
	ResultURL string `json:"result_url"`
	Error     string
}

func CleanURL(urlLink string) (string, error) {
	resp, err := http.PostForm(cleanURLAPI, url.Values{
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
	answer := &cleanURLAnswer{}
	err = json.Unmarshal(body, answer)
	if err != nil {
		return "", err
	}
	if answer.Error != "" {
		return "", errors.New(answer.Error)
	}
	return answer.ResultURL, nil
}

func Relink(urlLink string) (string, error) {
	jsn, err := json.Marshal(map[string]string{
		"url": urlLink,
	})
	if err != nil {
		return "", err
	}
	resp, err := http.Post(relinkAPI, "application/json", bytes.NewBuffer(jsn))
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

	cleanURL := "https://rel.ink/" + short
	return cleanURL, err
}
