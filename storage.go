package main

type Event struct{
    Owner string 
    Title string 
    Description string
    Tags []string
    Status string
    Prj string
}

type User struct{
   id string
   Name string
   Mail string
   pwd string
}


type Storage interface{
    initStorage() error
    //createAccount(*User) error
    //createEvent(*Event) error
    //deleteEvent(string) error
    //getAccountByID(string) (*Event, error)
}

