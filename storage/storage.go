package storage
import("eardis/types")


type Storage interface{
    CreateAccount(*types.User) (*types.User, error)
    GetEvents(userid string) ([]*types.Event, error)
    CreateEvent(e *types.Event) error
    DeleteEvent(ownerid string, eventid string) error
    PatchEvent(ownerid string,eventid string,e *types.Event) error
    Login(user *types.User)(string,error)
    GetUser(userid string)(*types.DisplayableUser,error)
    //updateUser(userid string)(*User,error) // (DisplayableUser,error)??
    GetNotifications(userid string)([]*types.Notifications,error)
    SearchUser(email string)(*types.DisplayableUser,error) 
    SendFriendRequestNotifications(notification types.Notifications)error 
    AcceptFriendRequest(notificationsid string, ownerid string)error 
    DeclineFriendRequest(notificationsid string,ownerid string)error 
    GetFriendsEvents(userid string) ([]*types.Event, error) 
   // createPrj(prj *Prj)error
   // deletePrj(prjid string)error
   // updatePrj(prj *Prj)error
   // addUserToPrj(prjuser PrjRole, prjid string)error
}
