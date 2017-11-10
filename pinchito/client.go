package pinchito

import (
	"encoding/json"
	"net/http"
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
