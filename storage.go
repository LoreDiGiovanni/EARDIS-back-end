package main

type DisplayableUser struct{
   ID string `json:"id" bson:"id,omitempty"`
   Username string `json:"username" bson:"username"`
   Email string `json:"email" bson:"email"`
}

type SearchUserRequest struct {
   Email string `json:"email" bson:"email"`
}

type Datetime struct{
    Date string `json:"date" bson:"date"`
    Time string `json:"time" bson:"time"`
}

// Enumerazione dei tipi di notifica
type NotificationType int

const (
	FriendRequest NotificationType = iota
	Comunication
)

type Notifications struct{
    ID string `json:"_id" bson:"_id,omitempty"`
    From string `json:"from" bson:"from"`
    To string `json:"to" bson:"to"`
    Type NotificationType `json:"type" bson:"type"`
    Description string `json:"description" bson:"description"`
}
type NotificationResponse struct{
    Notification_id string `json:"notification_id"`
    Response bool `json:"response"`
    Notification_type NotificationType `json:"notification_type"`
}

type PrjRole struct{
    Userid string `json:"userid" bson:"userid"`
    Admin bool `json:"admin" bson:"admin"`
    Edit_prj_name bool `json:"edit_prj_name" bson:"edit_prj_name"` 
    Edit_prj_description bool `json:"edit_prj_description" bson:"edit_prj_description"` 
    Edit_event_title bool `json:"edit_event_title" bson:"edit_event_title"`
    Edit_event_description bool `json:"edit_event_description" bson:"edit_event_description"`
    Edit_event_status bool `json:"change_event_status" bson:"change_event_status"`
    Delete_event bool `json:"delete_event" bson:"delete_event"`
    Create_event bool `json:"create_event" bson:"create_event"`
}

type Prj struct{
    ID string `json:"_id" bson:"_id,omitempty"`
    Name string `json:"prj_name" bson:"prj_name"`
    Description string  `json:"description" bson:"description"`
    StatusList []string `json:"status_list" bson:"status_list"`
    Members []PrjRole `json:"members" bson:"members"`
}

type Event struct{
    ID string `json:"id" bson:"_id,omitempty"`
    Owner string `json:"owner" bson:"owner"`
    Guests [] string `json:"guests" bson:"guests"`
    Title string  `json:"title" bson:"title"`
    Description string `json:"description" bson:"description"`
    Tags []string `json:"tags" bson:"tags"`
    Status string `json:"status" bson:"status"`
    Prjid string `json:"prj" bson:"prj"`
    Date Datetime `json:"date" bson:"date"`
    Deadline Datetime `json:"deadline" bson:"deadline"`
}

type User struct{
   ID string `json:"id" bson:"_id,omitempty"`
   Username string `json:"username" bson:"username"`
   Email string `json:"email" bson:"email"`
   PWD string `json:"pwd" bson:"pwd"`
   JWT string `json:"jwt" bson:"jwt"`
   Friends []string `json:"friends" bson:"friends"` 
   Prj_list []string `json:"prj_list" bson:"prj_list"`
}

type Storage interface{
    createAccount(*User) (*User, error)
    getEvents(userid string) ([]*Event, error)
    createEvent(e *Event) error
    deleteEvent(ownerid string, eventid string) error
    patchEvent(ownerid string,eventid string,e *Event) error
    login(user *User)(string,error)
    getUser(userid string)(*DisplayableUser,error)
    //updateUser(userid string)(*User,error) // (DisplayableUser,error)??
    getNotifications(userid string)([]*Notifications,error)
    searchUser(email string)(*DisplayableUser,error) 
    sendFriendRequestNotifications(notification Notifications)error 
    acceptFriendRequest(notificationsid string, ownerid string)error 
    declineFriendRequest(notificationsid string,ownerid string)error 
    getFriendsEvents(userid string) ([]*Event, error) 
   // createPrj(prj *Prj)error
   // deletePrj(prjid string)error
   // updatePrj(prj *Prj)error
   // addUserToPrj(prjuser PrjRole, prjid string)error
}
