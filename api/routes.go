package api

import (
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)

func (s *APIServer) SetupRoutes() *mux.Router {
    
    router := mux.NewRouter()
    handlers.CORS(
        handlers.AllowedOrigins([]string{"http://127.0.0.1:5173"}),//just for now
        handlers.AllowedMethods([]string{"GET","POST","DELETE","PATCH"}),
        handlers.AllowedHeaders([]string{"*"}),
    )(router)

    router.HandleFunc("/check", genericHandleFunc(s.HandleCheck))
    router.HandleFunc("/createAccount", genericHandleFunc(s.HandleCreateAccount))
    router.HandleFunc("/login", genericHandleFunc(s.HandleLogin))
    router.HandleFunc("/user", jwtHandleFunc(s.HandleUser))
    router.HandleFunc("/events", jwtHandleFunc(s.HandleEvents))
    router.HandleFunc("/events/id", jwtHandleFunc(s.HandleEventById))
    router.HandleFunc("/projects", jwtHandleFunc(s.HandleProjects))
    router.HandleFunc("/notifications", jwtHandleFunc(s.HandleNotifications))
    router.HandleFunc("/notifications/conferm", jwtHandleFunc(s.HandleNotificationsByID))
    router.HandleFunc("/friends", jwtHandleFunc(s.HandleFiriends))

    return router
}

