package pinchito

import (
	"encoding/json"
	"net/http"
	"net/url"
	"errors"
	"bytes"
	"os"
)

func Tapeta() (Log, error) {
	host := os.Getenv("PINCHITO_HOST")
	log := Log{}
	res, err := http.Get("http://" + host + "/json/random")
	if err != nil {
		return log, err
	}
	err = json.NewDecoder(res.Body).Decode(&log)
	return log, err
}

func Search(term string) (Log, error) {
	host := os.Getenv("PINCHITO_HOST")
	logs := []Log{}
	baseURL := "http://" + host + "/json/search?"
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

func UploadNewLog(uploadOp *JSONUploadOp) (int, error) {
	host := os.Getenv("PINCHITO_HOST")

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(uploadOp)
	res, err := http.Post("http://" + host + "/json/upload", "application/json", b)
	if err != nil {
		return -1, err
	}

	response := JSONUploadResult{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return -1, err
	}
	if !response.Result {
		return -1, errors.New("An error ocured while saving the log: " + response.ErrorMessage)
	}

	return response.IdPlog, nil

}

func GetUserFromPinchitoNick(nick string) (TgPinchitoUser, error) {
	for _, user := range tgPinchitoUsers {
		if user.PinNick == nick {
			return user, nil
		}
	}

	return TgPinchitoUser{}, errors.New("No User found with login " + nick)
}


func GetUserFromTelegramUsername(tgUsername string) (TgPinchitoUser, error) {
	for _, user := range tgPinchitoUsers {
		if user.TgUsername == tgUsername {
			return user, nil
		}
	}

	return TgPinchitoUser{}, errors.New("No User found with Telegram username " + tgUsername)
}

func GetPinchitoUsers() []TgPinchitoUser {
	return tgPinchitoUsers
}
