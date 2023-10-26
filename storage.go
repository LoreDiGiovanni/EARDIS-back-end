package main

type Event struct{
    Owner string `json:"owner" bson:"owner"`
    Title string  `json:"title" bson:"title"`
    Description string `json:"description" bson:"description"`
    Tags []string `json:"tags" bson:"tags"`
    Status string `json:"status" bson:"status"`
    Prj string `json:"prj" bson:"prj"`
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
    getEvents(string) ([]*Event, error)
    createEvent(*Event) error
    //deleteEvent(string) error
}

