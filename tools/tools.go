package tools 

import(
    "net/http"
    "encoding/json"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
    w.Header().Add("Content-Type","application/json")
    w.WriteHeader(status)
    if(v != nil){
        return json.NewEncoder(w).Encode(v);
    }else{
        return nil 
    }
}


