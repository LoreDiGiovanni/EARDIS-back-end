package main

import(
    "net/http"
    "encoding/json"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
    w.Header().Add("Content-Type","application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(v);
}


