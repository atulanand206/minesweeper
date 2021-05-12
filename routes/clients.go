package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/atulanand206/minesweeper/objects"
)

const (
	acceptKey       = "Accept"
	applicationJson = "application/json"
)

func GetUsers(usernames []string) ([]objects.User, error) {
	client := &http.Client{}
	hostname := os.Getenv("USERS_URL")
	endpoint := "/users"
	url := "http://" + hostname + endpoint
	requestByte, _ := json.Marshal(usernames)
	requestReader := bytes.NewReader(requestByte)
	request, err := http.NewRequest("GET", url, requestReader)
	var ob []objects.User
	if err != nil {
		return ob, err
	}
	request.Header.Add(acceptKey, applicationJson)
	request.Header.Add(contentTypeKey, applicationJson)
	response, err := client.Do(request)
	fmt.Println(response)
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
