package routes

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/atulanand206/minesweeper/objects"
)

const (
	acceptKey       = "Accept"
	applicationJson = "application/json"
)

func GetUser(username string) (objects.User, error) {
	client := &http.Client{}
	hostname := os.Getenv("USERS_URL")
	endpoint := os.Getenv("USERS_PATH_GET_USER")
	url := "http://" + hostname + endpoint + username
	request, err := http.NewRequest("GET", url, nil)
	var ob objects.User
	if err != nil {
		return ob, err
	}
	request.Header.Add(acceptKey, applicationJson)
	request.Header.Add(contentTypeKey, applicationJson)
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
