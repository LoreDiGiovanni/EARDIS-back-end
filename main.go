package main

import (
	"log"
    "eardis/api"
    "eardis/storage"
)

func main(){
    store, err := storage.NewMongoStore()
    if err!=nil{
        log.Fatal("[X] Storeg connection error")
    }else{
        server := api.NewAPIServer(":3000",store)
        log.Println("[v] Back end started")
        server.Run()
    }
}
