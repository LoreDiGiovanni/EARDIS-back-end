package main

type Event struct{
    Owner string `json:"owner"`
    Title string  `json:"title"`
    Description string `json:"Description"`
    Tags []string `json:"tags"`
    Status string `json:"status"`
    Prj string `json:"prj"`
}

type User struct{
   ID string `json:"id" bson:"id,omitempty"`
   Username string `json:"username" bson:"username"`
   Email string `json:"email" bson:"email"`
   PWD string `json:"pwd" bson:"pwd"`
   JWT string `json:"jwt" bson:"jwt"`
}

type Storage interface{
    createAccount(*User) (*User, error)
    getEvent(*User) (*Event, error)
    //createEvent(*Event) error
    //deleteEvent(string) error
}

