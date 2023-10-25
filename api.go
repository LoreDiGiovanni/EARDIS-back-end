package main

import (
	"fmt"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct{
    address string
    store Storage
}
type ApiError struct {
    Error string
}

type genericHandle func(http.ResponseWriter,*http.Request) error

func genericHandleFunc(f genericHandle) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        err := f(w,r)
        if err != nil{
            WriteJSON(w,http.StatusBadRequest,ApiError{Error: err.Error()})
        }
    }
}

func NewAPIServer(address string,store Storage) *APIServer{
    return &APIServer{address: address, store: store}
}

func (s* APIServer) Run() {
    router := mux.NewRouter() 
    router.HandleFunc("/createAccount",genericHandleFunc(s.HandleCreateAccount))
    router.HandleFunc("/events",jwtHandleFunc(s.HandleEvents))
    http.ListenAndServe(s.address,router)
}

func (s* APIServer) HandleEvents(w http.ResponseWriter, r *http.Request,t *jwt.Token) error{
    switch r.Method {
        case "GET": return s.getEvents(w,r,t) 
        case "POST": return s.createEvent(w,r,t) 
        case "PATCH": return s.patchEvent(w,r,t) 
        case "DELETE": return s.deleteEvent(w,r,t) 
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s* APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error{
    if r.Method == "POST"{
        return s.createAccount(w,r) 
    }else{
        return fmt.Errorf("Method not allowed %s", r.Method)
    }
}

func (s* APIServer) createEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}
func (s* APIServer) patchEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}
func (s* APIServer) deleteEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}
func (s* APIServer) getEvents(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}

func (s* APIServer) createAccount(w http.ResponseWriter,r *http.Request) error{
    username := r.Header.Get("username")
    email := r.Header.Get("email")
    pwd := r.Header.Get("pwd")
    user := User{Username: username,Email: email,PWD: pwd}
    newuser,err := s.store.createAccount(&user)
    if err!=nil{
        return WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Email or Username already used!"})
    }else{
        token := struct{Token string `json:"token"`}{Token: newuser.JWT}
        return WriteJSON(w,http.StatusOK,token)
    }
}
