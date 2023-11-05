package authentication

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"errors"
	"time"
	"fmt"
	"strconv"
)

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginJsonResponse struct {
	Error   bool     `json:"error"`
	Message string   `json:"message"`
	User    User     `json:"user,omitempty"`
}

func Login(entry LoginRequest) (error, User) {
	//creare a variable we'll read response.Body into
	var jsonFromService LoginJsonResponse
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	// call the service
	request, err := http.NewRequest("POST", "http://localhost:5442/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Couldn`t login, something wrong with auth service")
		return err, jsonFromService.User
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err, jsonFromService.User
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("Wrong login or password, code: %s", strconv.Itoa(response.StatusCode))), jsonFromService.User
	}

	//decode the json
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		return err, jsonFromService.User
	}
	if jsonFromService.Error == true {
		return errors.New(fmt.Sprintf("Wrong login or password, code2: %s", strconv.Itoa(response.StatusCode))), jsonFromService.User
	}

	return nil, jsonFromService.User
}

func CheckRedisSession() error {
	jsonData, _ := json.MarshalIndent("", "", "\t")
	request, err := http.NewRequest("GET", "http://localhost:4111/auth/test", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Could not check customer session")
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusOK {
		return errors.New("Customer session is not initialized")
	}
	return nil
}