package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"eardis/storage"
	"eardis/types"
	"eardis/tools"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct{
    address string
    store storage.Storage
}
type ApiError struct {
    Error string
}

type genericHandle func(http.ResponseWriter,*http.Request) error

func genericHandleFunc(f genericHandle) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        err := f(w,r)
        if err != nil{
            tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: err.Error()})
        }
    }
}

func NewAPIServer(address string,store storage.Storage) *APIServer{
    return &APIServer{address: address, store: store}
}

func (s* APIServer) Run() {
    router := mux.NewRouter() 
    router.HandleFunc("/createAccount",genericHandleFunc(s.HandleCreateAccount))
    router.HandleFunc("/login",genericHandleFunc(s.HandleLogin))
    router.HandleFunc("/user",jwtHandleFunc(s.HandleUser))
    router.HandleFunc("/events",jwtHandleFunc(s.HandleEvents))
    router.HandleFunc("/events/{eventid}",jwtHandleFunc(s.HandleEventById))
    router.HandleFunc("/projects",jwtHandleFunc(s.HandleProjects))
    router.HandleFunc("/notifications",jwtHandleFunc(s.HandleNotifications))
    router.HandleFunc("/notifications/conferm",jwtHandleFunc(s.HandleNotificationsByID))
    router.HandleFunc("/friends",jwtHandleFunc(s.HandleFiriends))
    http.ListenAndServe(s.address,router)
}

func (s* APIServer) HandleFiriends(w http.ResponseWriter, r *http.Request,t *jwt.Token) error{
    switch r.Method {
        //case "GET": return s.getUser(w,r,t) 
        //case "POST": return s.searchUser(w,r,t) 
        
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}
func (s* APIServer) HandleNotificationsByID(w http.ResponseWriter, r *http.Request,t *jwt.Token) error{
    switch r.Method {
        case "POST": return s.replyToNotification(w,r,t) 
        
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s* APIServer) replyToNotification(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    claims := t.Claims.(jwt.MapClaims)
    userid := claims["id"].(string)
    var nresponse types.NotificationResponse
    err := json.NewDecoder(r.Body).Decode(&nresponse);if err != nil {return err}
    defer r.Body.Close()
    switch nresponse.Notification_type{
        case types.FriendRequest: {
            if nresponse.Response{
                err := s.store.AcceptFriendRequest(nresponse.Notification_id,userid); if err!= nil{
                    return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Notification not exist"})
                }else{
                    return tools.WriteJSON(w,http.StatusOK,nil)
                }
            }else{
                err := s.store.DeclineFriendRequest(nresponse.Notification_id,userid); if err!= nil{
                    return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Notification not exist"})
                }else{
                    return tools.WriteJSON(w,http.StatusOK,nil)
                }

            }
            
        }
    }
    return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Invalid Notification type"})
}

func (s* APIServer) HandleNotifications(w http.ResponseWriter, r *http.Request,t *jwt.Token) error{
    switch r.Method {
        case "GET": return s.getNotifications(w,r,t) 
        case "POST": return s.sendNotifications(w,r,t) 
        
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s* APIServer) getNotifications(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
     claims := t.Claims.(jwt.MapClaims)
     userid := claims["id"].(string)
     notifications,err := s.store.GetNotifications(userid); if err != nil{
        return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "User not found"})
     }else{
         if len(notifications)<=0{
            return tools.WriteJSON(w,http.StatusOK,nil)
         }else{
            return tools.WriteJSON(w,http.StatusOK,notifications)
         }
     }
}

func (s* APIServer) sendNotifications(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    claims := t.Claims.(jwt.MapClaims)
    userid := claims["id"].(string)
    var message types.Notifications 
    err := json.NewDecoder(r.Body).Decode(&message);if err != nil {return err}
    defer r.Body.Close()
    message.From = userid
    switch message.Type{
        case types.FriendRequest: {
            err := s.store.SendFriendRequestNotifications(message); if err!= nil{
                return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "User does not exist"})
            }else{
                return tools.WriteJSON(w,http.StatusOK,nil)
            }
        }
    }
    return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Invalid Notification type"})
}

func (s* APIServer) HandleProjects(w http.ResponseWriter, r *http.Request,t *jwt.Token) error{
    switch r.Method {
        //case "GET": return s.getUser(w,r,t) 
        //case "POST": return s.searchUser(w,r,t) 
        
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s* APIServer) HandleUser(w http.ResponseWriter, r *http.Request,t *jwt.Token) error{
    switch r.Method {
        case "GET": return s.getUser(w,r,t) 
        case "POST": return s.searchUser(w,r,t) 
        
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s* APIServer) getUser(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
     claims := t.Claims.(jwt.MapClaims)
     userid := claims["id"].(string)
     user,err := s.store.GetUser(userid); if err != nil{
        log.Println(err)
        return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Account does not exist"})
    }else{
        return tools.WriteJSON(w,http.StatusOK,user)
    }
}

func (s* APIServer) searchUser(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    var search types.SearchUserRequest
    err := json.NewDecoder(r.Body).Decode(&search);if err != nil {return err}
    defer r.Body.Close()
    user,err := s.store.SearchUser(search.Email);if err!= nil{
        return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Account does not exist"})
    }else{
        return tools.WriteJSON(w,http.StatusOK,user)
    }
}

func (s* APIServer) HandleEvents(w http.ResponseWriter, r *http.Request,t *jwt.Token) error{
    switch r.Method {
        case "GET": return s.getEvents(w,r,t) 
        case "POST": return s.createEvent(w,r,t) 
        //case "PATCH": return s.patchEvent(w,r,t) 
        //case "DELETE": return s.deleteEvent(w,r,t) 
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s* APIServer) getEvents(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    var events []*types.Event 
    var err error
    claims := t.Claims.(jwt.MapClaims)
    id := claims["id"].(string)
    events, err =  s.store.GetEvents(id); if err!= nil{
        log.Println("Function: getEvents, id: ",id,", Error: ",err)
        return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "User not found"})
    }else{
        return tools.WriteJSON(w,http.StatusOK,events)
    }
}

func (s* APIServer) createEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    claims := t.Claims.(jwt.MapClaims)
    ownerid := claims["id"].(string)
    var event types.Event = types.Event{Owner: ownerid}
    err := json.NewDecoder(r.Body).Decode(&event);if err != nil {return err}
    defer r.Body.Close()
    s.store.CreateEvent(&event)
    return tools.WriteJSON(w,http.StatusOK,nil)
}


func (s* APIServer) HandleEventById(w http.ResponseWriter, r *http.Request,t *jwt.Token) error{
    switch r.Method {
        //case "GET": return s.getEvent(w,r,t) 
        case "PATCH": return s.patchEvent(w,r,t) 
        case "DELETE": return s.deleteEvent(w,r,t) 
    }
    return fmt.Errorf("Method not allowed %s", r.Method)
}
func (s* APIServer) patchEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    claims := t.Claims.(jwt.MapClaims)
    ownerid := claims["id"].(string)
    eventid := mux.Vars(r)["eventid"]
    var event types.Event 
    err := json.NewDecoder(r.Body).Decode(&event);if err != nil {return err}
    defer r.Body.Close()
    event.Owner = ownerid
    event.ID = eventid
    err = s.store.PatchEvent(ownerid,eventid,&event); if err != nil{
        return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "non-existent event, impossible to update"})
    }else{
        return tools.WriteJSON(w,http.StatusOK,nil)
    }
}

func (s* APIServer) deleteEvent(w http.ResponseWriter,r *http.Request,t *jwt.Token) error{
    claims := t.Claims.(jwt.MapClaims)
    ownerid := claims["id"].(string)
    eventid := mux.Vars(r)["eventid"]
    err :=  s.store.DeleteEvent(ownerid,eventid); if err!= nil{
        log.Println("Function: deleteEvent ","Error: ",err)
        return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Impossible to delete the event"})
    }else{
        return tools.WriteJSON(w,http.StatusOK,nil)
    }
}

func (s* APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error{
    if r.Method == "POST"{
        return s.createAccount(w,r) 
    }else{
        return fmt.Errorf("Method not allowed %s", r.Method)
    }
}

func (s* APIServer) createAccount(w http.ResponseWriter,r *http.Request) error{
    var user types.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {return err}
    defer r.Body.Close()
    newuser,err := s.store.CreateAccount(&user); if err!=nil{
        return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "Email or Username already used!"})
    }else{
        token := struct{Token string `json:"token"`}{Token: newuser.JWT}
        return tools.WriteJSON(w,http.StatusOK,token)
    }
}

func (s* APIServer) HandleLogin(w http.ResponseWriter, r *http.Request) error{
    if r.Method == "POST"{
        return s.login(w,r) 
    }else{
        return fmt.Errorf("Method not allowed %s", r.Method)
    }
}

func (s* APIServer) login(w http.ResponseWriter,r *http.Request) error{
    var user types.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {return err}
    defer r.Body.Close()
    token,err := s.store.Login(&user);if err !=nil{
        return tools.WriteJSON(w,http.StatusBadRequest,ApiError{Error: "User formatting error"})
    }else{
        token := struct{Token string `json:"token"`}{Token: token}
        return tools.WriteJSON(w,http.StatusOK,token)
    }
}






