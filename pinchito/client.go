package pinchito

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func Tapeta() (Log, error) {
	log := Log{}
	res, err := http.Get("http://go.pinchito.com/json/random")
	if err != nil {
		return log, err
	}
	err = json.NewDecoder(res.Body).Decode(&log)
	return log, err
}

func Search(term string) (Log, error) {
	logs := []Log{}
	baseURL := "http://go.pinchito.com/json/search?"
	params := url.Values{}
	params.Add("s", term)

	finalURL := baseURL + params.Encode()
	res, err := http.Get(finalURL)
	if err != nil {
		return Log{}, err
	}
	err = json.NewDecoder(res.Body).Decode(&logs)
	if err != nil {
		return Log{}, err
	}
	if len(logs) == 0 {
		return Log{}, nil
	}
	return logs[0], nil
}
