package pinchito

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	PinchitoHost string
}

func (c *Client) Tapeta() (Log, error) {
	log := Log{}
	res, err := http.Get("http://" + c.PinchitoHost + "/json/random")
	if err != nil {
		return log, err
	}
	err = json.NewDecoder(res.Body).Decode(&log)
	return log, err
}

func (c *Client) Search(term string) (Log, error) {
	logs := []Log{}
	baseURL := "http://" + c.PinchitoHost + "/json/search?"
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

func (c *Client) UploadNewLog(uploadOp *JSONUploadOp) (string, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(uploadOp)
	res, err := http.Post("http://"+c.PinchitoHost+"/json/upload", "application/json", b)
	if err != nil {
		return "", err
	}

	response := JSONUploadResult{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	if !response.Result {
		return "", errors.New("An error ocured while saving the log: " + response.ErrorMessage)
	}

	return "http://" + c.PinchitoHost + "/" + strconv.Itoa(response.IdPlog), nil

}

func (c *Client) GetUserFromPinchitoNick(nick string) (TgPinchitoUser, error) {
	for _, user := range tgPinchitoUsers {
		if user.PinNick == nick {
			return user, nil
		}
	}

	return TgPinchitoUser{}, errors.New("No User found with login " + nick)
}

func (c *Client) GetUserFromTelegramUsername(tgUsername string) (TgPinchitoUser, error) {
	for _, user := range tgPinchitoUsers {
		if user.TgUsername == tgUsername {
			return user, nil
		}
	}

	return TgPinchitoUser{}, errors.New("No User found with Telegram username " + tgUsername)
}

func (c *Client) GetPinchitoUsers() []TgPinchitoUser {
	return tgPinchitoUsers
}
