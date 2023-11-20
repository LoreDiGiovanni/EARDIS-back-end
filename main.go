package main

import (
	"log"
    "eardis/api"
    "eardis/storage"
)

func main(){
    log.Println("[v] Back end started")
    store, err := storage.NewMongoStore()
    if err!=nil{
        log.Fatal("[X] Storeg connection error")
    }
    server := api.NewAPIServer(":3000",store)
    server.Run()
}
