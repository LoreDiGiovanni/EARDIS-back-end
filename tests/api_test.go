package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
    "eardis/types"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite
    Token string
    Userid string
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APISuite))
}


func (s *APISuite) Test1CheckEndpointGET() {
	req, err := http.NewRequest("GET", "http://localhost:3000/check", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

    s.Equal(http.StatusOK,resp.StatusCode)
}


func (s *APISuite) Test2CreateAccountEndpointPOST() {
    data := map[string]string{
		"username": "test_user",
		"email": "test_user@gmail.com",
		"pwd": "test_pwd",
	}
	jsonData, err := json.Marshal(data);if err != nil {
		s.T().Fatal(err)
	}
    req, err := http.NewRequest("POST", "http://localhost:3000/createAccount",bytes.NewBuffer(jsonData))
	if err != nil {
		s.T().Fatal(err)
	}

	// Esegue la richiesta
	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		s.T().Fatal(err)
	}

    var tokenResponse types.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse);if err != nil {
		s.T().Fatal(err)
	}

	// Verifica lo stato della risposta
    s.Equal(http.StatusOK,resp.StatusCode)
    s.NotEqual("",tokenResponse.Token)
    s.Token = tokenResponse.Token
}

func (s *APISuite) Test3UserEndpoint1Get() {
    req, err := http.NewRequest("GET", "http://localhost:3000/user",nil)
	if err != nil {
		s.T().Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", s.Token)
	// Esegue la richiesta
	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		s.T().Fatal(err)
	}

    var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user);if err != nil {
		s.T().Fatal(err)
	}

	// Verifica lo stato della risposta
    s.Equal(http.StatusOK,resp.StatusCode)
    s.Equal("test_user",user.Username)
    s.Equal("test_user@gmail.com",user.Email)
    s.Userid = user.ID
}

func (s *APISuite) Test4UserEndpointDelete() {
    req, err := http.NewRequest("DELETE", "http://localhost:3000/user",nil)
	if err != nil {
		s.T().Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", s.Token)

	// Esegue la richiesta
	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		s.T().Fatal(err)
	}
    s.Equal(http.StatusOK,resp.StatusCode)
}
