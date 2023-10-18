package main

import (
    "fmt"
)

func main(){
    fmt.Println("[v] Back end started")
    server := NewAPIServer(":3000")
    server.Run()

}
