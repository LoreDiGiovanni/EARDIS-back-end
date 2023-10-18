package main

import (
	"fmt"
	"net/http"

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
func makeHTTPHandleFunc(f genericHandle) http.HandlerFunc {
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
    router.HandleFunc("/events",makeHTTPHandleFunc(s.HandleEvents))
    http.ListenAndServe(s.address,router)
}

func (s* APIServer) HandleEvents(w http.ResponseWriter, r *http.Request) error{
    switch r.Method {
        case "POST": return createEvent(w) 
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}

func createEvent(w http.ResponseWriter) error{
    return WriteJSON(w,http.StatusOK,Event{Title: "Test"})
}
