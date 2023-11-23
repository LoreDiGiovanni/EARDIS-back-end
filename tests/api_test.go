package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
    "eardis/types"
    "eardis/api"
    "eardis/storage"
)

var Token string
var Token2 string
var Userid string
var Userid2 string
var Event types.Event 
var notificationId string
var notificationType int

func TestAPI(t *testing.T) {
    store, err := storage.NewMongoStore()
    if err!=nil{
        panic(nil)
    }
    server := api.NewAPIServer(":3000",store)
    go server.Run()
}


func TestCheckEndpointGET(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/check", nil)
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK != resp.StatusCode{
        t.Fatal()
    }
}


func TestCreateAccount(t *testing.T) {
    data := map[string]string{
		"username": "test_user",
		"email": "test_user@gmail.com",
		"pwd": "test_pwd",
	}
	jsonData, err := json.Marshal(data);if err != nil {
		t.Fatal(err)
	}
    req, err := http.NewRequest("POST", "http://localhost:3000/createAccount",bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// Esegue la richiesta
	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    var tokenResponse types.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK != resp.StatusCode{
        t.Fatal()
    }else if "" == tokenResponse.Token{
        t.Fatal()
    }
    Token = tokenResponse.Token
}

func  TestGetAccount(t *testing.T){
    req, err := http.NewRequest("GET", "http://localhost:3000/user",nil)
	if err != nil {
		t.Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token)
	// Esegue la richiesta
	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user);if err != nil {
		t.Fatal(err)
	}

	// Verifica lo stato della risposta
    if http.StatusOK != resp.StatusCode{
		t.Fatal()
    }else if "test_user"!= user.Username{
		t.Fatal()
    }else if "test_user@gmail.com"!=user.Email{
		t.Fatal()
    }
    Userid = user.ID
}

func TestCreateSecondUser(t *testing.T) {
    data := map[string]string{
		"username": "test_user2",
		"email": "test_user2@gmail.com",
		"pwd": "test_pwd2",
	}
	jsonData, err := json.Marshal(data);if err != nil {
		t.Fatal(err)
	}
    req, err := http.NewRequest("POST", "http://localhost:3000/createAccount",bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    var tokenResponse types.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK != resp.StatusCode{
		t.Fatal(err)
    }else if "" == tokenResponse.Token{
		t.Fatal(err)
    }
    Token2 = tokenResponse.Token
}

func TestSearchUser(t *testing.T) {
    data := map[string]string{
		"email": "test_user2@gmail.com",
	}
	jsonData, err := json.Marshal(data);if err != nil {
		t.Fatal(err)
	}

    req, err := http.NewRequest("POST", "http://localhost:3000/user",bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    var user types.DisplayableUser
	err = json.NewDecoder(resp.Body).Decode(&user);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK != resp.StatusCode{
		t.Fatal(err)
    }else if "test_user2" != user.Username{
		t.Fatal(err)
    }else if "test_user2@gmail.com" != user.Email{
		t.Fatal(err)
    }
    Userid2 = user.ID
}


func TestLogin(t *testing.T) {
    data := map[string]string{
		"email": "test_user@gmail.com",
		"pwd": "test_pwd",
	}
	jsonData, err := json.Marshal(data);if err != nil {
		t.Fatal(err)
	}

    req, err := http.NewRequest("POST", "http://localhost:3000/login",bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    var tokenResponse types.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK!=resp.StatusCode{
		t.Fatal(err)
    }else if Token != tokenResponse.Token {
		t.Fatal(err)
    }
}

func TestCreateEvent(t *testing.T){
    data := map[string]string{
		"title": "test_title",
	}

	jsonData, err := json.Marshal(&data);if err != nil {
		t.Fatal(err)
	}

    req, err := http.NewRequest("POST", "http://localhost:3000/events",bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK!=resp.StatusCode{
		t.Fatal(err)
    }
}

func TestGetEvents(t *testing.T) {
    req, err := http.NewRequest("GET", "http://localhost:3000/events",nil)
	if err != nil {
		t.Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    var events []types.Event
    err = json.NewDecoder(resp.Body).Decode(&events);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK!=resp.StatusCode{
		t.Fatal(err)
    }else if len(events)<=0{
		t.Fatal(err)
    }else if "test_title" != events[0].Title {
		t.Fatal(err)
    }
    Event = events[0]
}

func TestUpdateEvent(t *testing.T) {
    Event.Description = "new description"
    jsonData, err := json.Marshal(&Event);if err != nil {
		t.Fatal(err)
	}
    req, err := http.NewRequest("PATCH", 
                                "http://localhost:3000/events/id",
                                 bytes.NewBuffer(jsonData))

	if err != nil {
		t.Fatal(err)
	}
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK!=resp.StatusCode{
		t.Fatal(err)
    }
}




func TestDeleteEvent(t *testing.T) {
    data := map[string]string{
		"id": Event.ID,
	}

    jsonData, err := json.Marshal(&data);if err != nil {
		t.Fatal(err)
	}

    req, err := http.NewRequest("DELETE", 
                                "http://localhost:3000/events/id",
                                 bytes.NewBuffer(jsonData))

	if err != nil {
		t.Fatal(err)
	}
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK != resp.StatusCode{
        t.Fatal()
    }
}

func  TestSendFriendRequeste(t *testing.T){
    data := map[string]interface{}{
        "to":   Userid2,
        "type": 0,
    }

	jsonData, err := json.Marshal(&data);if err != nil {
		t.Fatal(err)
	}

    req, err := http.NewRequest("POST", "http://localhost:3000/notifications",bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token)


	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK != resp.StatusCode{
		t.Fatal()
    }
}

func  TestGetNotifications(t *testing.T){
   	req, err := http.NewRequest("GET", "http://localhost:3000//notifications",nil)
	if err != nil {
		t.Fatal(err)
	}
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token2)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    var notifications []types.Notifications
	err = json.NewDecoder(resp.Body).Decode(&notifications);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK != resp.StatusCode{
		t.Fatal()
    }else if len(notifications)<=0{
		t.Fatal()
    }
    notificationId = notifications[0].ID
    notificationType = int(notifications[0].Type)
}

func  TestAcceptFriendRequest(t *testing.T){
    data := map[string]interface{}{
        "notification_id":   notificationId,
        "response": true,
        "notification_type": notificationType,
    }

	jsonData, err := json.Marshal(&data);if err != nil {
		t.Fatal(err)
	}

   	req, err := http.NewRequest("POST", "http://localhost:3000/notifications/conferm",bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token2)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}

    if http.StatusOK != resp.StatusCode{
		t.Fatal()
    }
}

func TestDeleteSecondUser(t *testing.T) {
    req, err := http.NewRequest("DELETE", "http://localhost:3000/user",nil)
	if err != nil {
		t.Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token2)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}
    if http.StatusOK != resp.StatusCode{
        t.Fatal()
    }
}

func TestDeleteUser(t *testing.T) {
    req, err := http.NewRequest("DELETE", "http://localhost:3000/user",nil)
	if err != nil {
		t.Fatal(err)
	}

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("x-jwt-token", Token)

	client := http.Client{}
	resp, err := client.Do(req);if err != nil {
		t.Fatal(err)
	}
    if http.StatusOK != resp.StatusCode{
        t.Fatal()
    }
}

