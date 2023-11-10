package main

type Datetime struct{
    Date string `json:"date" bson:"date"`
    Time string `json:"time" bson:"time"`
}
type Prj struct{
    PrjName string `json:"prj_name" bson:"prj_name"`
    StatusList []string `json:"status_list" bson:"status_list"`
}

type Event struct{
    Owner string `json:"owner" bson:"owner"`
    Guests [] string `json:"guests" bson:"guests"`
    Title string  `json:"title" bson:"title"`
    Description string `json:"description" bson:"description"`
    Tags []string `json:"tags" bson:"tags"`
    Status string `json:"status" bson:"status"`
    Prj Prj `json:"prj" bson:"prj"`
    Date Datetime `json:"date" bson:"date"`
    Deadline Datetime `json:"deadline" bson:"deadline"`
}

type User struct{
   ID string `json:"id" bson:"id,omitempty"`
   Username string `json:"username" bson:"username"`
   Email string `json:"email" bson:"email"`
   PWD string `json:"pwd" bson:"pwd"`
   JWT string `json:"jwt" bson:"jwt"`
   Friends []string `json:"friends" bson:"friends"` 
}

type Storage interface{
    createAccount(*User) (*User, error)
    getEvents(string) ([]*Event, error)
    createEvent(*Event) error
    //deleteEvent(string) error
}

