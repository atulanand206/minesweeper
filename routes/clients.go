package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	net "github.com/atulanand206/go-network"
	"github.com/atulanand206/minesweeper/objects"
)

func GetUsers(usernames []string, headers http.Header) ([]objects.User, error) {
	client := &http.Client{}
	hostname := os.Getenv("USERS_URL")
	endpoint := "/users"
	url := "http://" + hostname + endpoint
	requestByte, _ := json.Marshal(usernames)
	requestReader := bytes.NewReader(requestByte)
	request, err := http.NewRequest(http.MethodGet, url, requestReader)
	for x, v := range headers {
		for _, y := range v {
			request.Header.Add(x, y)
		}
	}
	var ob []objects.User
	if err != nil {
		return ob, err
	}
	request.Header.Add(net.Accept, net.ApplicationJson)
	request.Header.Add(net.ContentTypeKey, net.ApplicationJson)
	response, err := client.Do(request)
	if err != nil {
		return ob, err
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&ob)
	if err != nil {
		return ob, err
	}
	return ob, nil
}
