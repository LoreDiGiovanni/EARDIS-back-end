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
        case "GET": return getEvent(w,r,t) 
        case "POST": return createEvent(w,r,t) 
        case "PATCH": return patchEvent(w,r,t) 
        case "DELETE": return deleteEvent(w,r,t) 
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s* APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error{
    if r.Method == "POST"{
        return createAccount(w,r) 
    }else{
        return fmt.Errorf("Method not allowed %s", r.Method)
    }
}

func createEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}
func patchEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}
func deleteEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}
func getEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}

func createAccount(w http.ResponseWriter,r *http.Request) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}
