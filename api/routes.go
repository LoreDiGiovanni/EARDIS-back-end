package api

import (
    "github.com/gorilla/mux"
)

func (s *APIServer) SetupRoutes() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/check", genericHandleFunc(s.HandleCheck))
    router.HandleFunc("/createAccount", genericHandleFunc(s.HandleCreateAccount))
    router.HandleFunc("/login", genericHandleFunc(s.HandleLogin))
    router.HandleFunc("/user", jwtHandleFunc(s.HandleUser))
    router.HandleFunc("/events", jwtHandleFunc(s.HandleEvents))
    router.HandleFunc("/events/{eventid}", jwtHandleFunc(s.HandleEventById))
    router.HandleFunc("/projects", jwtHandleFunc(s.HandleProjects))
    router.HandleFunc("/notifications", jwtHandleFunc(s.HandleNotifications))
    router.HandleFunc("/notifications/conferm", jwtHandleFunc(s.HandleNotificationsByID))
    router.HandleFunc("/friends", jwtHandleFunc(s.HandleFiriends))

    return router
}

