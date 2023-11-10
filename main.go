package main

import (
	"log"
)

func main(){
    log.Println("[v] Back end started")
    store, err := newMongoStore()
    if err!=nil{
        log.Fatal("[X] Storeg connection error")
    }
    server := NewAPIServer(":3000",store)
    server.Run()
}
