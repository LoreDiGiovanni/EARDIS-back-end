package storage

import (
	"context"
	"eardis/tools"
	"eardis/types"
	"errors"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type mongoStore struct{
    db *mongo.Client
}

func NewMongoStore() (*mongoStore,error,){
    uri := os.Getenv("MONGODB_URI")
    client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
    if err!= nil {
        return nil,err
    }else{
        return &mongoStore{db: client},nil
    }
}

func (s *mongoStore) CreateAccount(user *types.User) (*types.User,error){
    var err error
    coll := s.db.Database("eardis").Collection("users")
    res,err := coll.InsertOne(context.TODO(),user)
    if err != nil{
        return nil,err
    }else{
        id := res.InsertedID.(primitive.ObjectID)
        user.ID = id.Hex()
        user.JWT,err= tools.CreateUserJWT(user)
        filter := bson.D{{"_id", id}}
        update := bson.D{{"$set", bson.D{{"jwt", user.JWT}}}}
        _, err := coll.UpdateOne(context.TODO(), filter, update)
        if err!=nil {
            return nil,err;
        }
        log.Println("[V] New user: ",user.ID)
        return user,nil
    }
}

func (s *mongoStore) CreateEvent(e *types.Event) error{
    coll := s.db.Database("eardis").Collection("events")
    res,err := coll.InsertOne(context.TODO(),e)
    if err != nil{
        log.Println(err)
        return err
    }else{
        log.Println("[V] New Event: ",res.InsertedID.(primitive.ObjectID))
        return nil
    }
}

func (s *mongoStore) GetEvents(userid string) ([]*types.Event, error){
    coll := s.db.Database("eardis").Collection("events")
    filter := bson.D{{"owner", userid}}
    cursor, err := coll.Find(context.TODO(), filter); if err!= nil{
        return nil,err    
    }else{
        var results []*types.Event
        if err = cursor.All(context.TODO(), &results); err != nil {
            return nil,err
	    }else{
            fireds_resout,err := s.GetFriendsEvents(userid);if err != nil{
                return results,nil
            }else{
                results := append(results, fireds_resout...)
                return results,nil
            }
        }
    }
}
func (s *mongoStore) DeleteEvent(ownerid string ,eventid string) error{
    coll := s.db.Database("eardis").Collection("events")
    objectID, err := primitive.ObjectIDFromHex(eventid); if err != nil {return err}
    filter := bson.D{{"_id", objectID},{"owner", ownerid}}
    _, err = coll.DeleteOne(context.TODO(), filter); if err != nil {
        return err
    }else{
        return nil
    }
}

func (s *mongoStore) PatchEvent(ownerid string ,eventid string,e *types.Event) error{
    coll := s.db.Database("eardis").Collection("events")
    objectID, err := primitive.ObjectIDFromHex(eventid); if err != nil {
		return errors.New("Invalid event id") 
	}
    e.ID = ""
    filter := bson.D{{"_id", objectID},{"owner", ownerid}}
    update := bson.D{{"$set", e}}
    _ , err = coll.UpdateOne(context.TODO(), filter, update); if err != nil {
        return err
    }else{
        return nil
    }
}

func (s *mongoStore)Login(user *types.User) (string,error){
    coll := s.db.Database("eardis").Collection("users")
    filter := bson.D{{"email",user.Email}}
    var dbuser types.User
    err := coll.FindOne(context.TODO(),filter).Decode(&dbuser); if err != nil{
        return "", err 
    }else{
        if dbuser.PWD == tools.RiGeneratePwd(user.PWD,dbuser.Salt){
            return dbuser.JWT, nil
        }else{
            return "", fmt.Errorf("Wrong email or password")
        }
    }
}

func (s *mongoStore)SearchUser(email string) (*types.DisplayableUser,error){
    coll := s.db.Database("eardis").Collection("users")
    filter := bson.D{{"email",email}}
    var friend types.User
    err := coll.FindOne(context.TODO(),filter).Decode(&friend); if err != nil{
        return nil,err 
    }else{
        var duser types.DisplayableUser = types.DisplayableUser{ID: friend.ID,Username: friend.Username,Email: friend.Email}
        return &duser, nil
    }
}

func (s *mongoStore)GetUser(userid string)(*types.DisplayableUser,error){
    coll := s.db.Database("eardis").Collection("users")
    objectID, err := primitive.ObjectIDFromHex(userid); if err != nil {
		return nil,errors.New("Invalid event id") 
	}
    filter := bson.D{{"_id",objectID}}
    var user types.User
    err = coll.FindOne(context.TODO(),filter).Decode(&user); if err != nil{
        return nil,err 
    }else{
        var duser types.DisplayableUser = types.DisplayableUser{ID: user.ID,Username: user.Username,Email: user.Email}
        return &duser, nil
    }
}

func (s *mongoStore)GetNotifications(userid string)([]*types.Notifications,error){    
    coll := s.db.Database("eardis").Collection("notifications")
    filter := bson.D{{"to", userid}}
    cursor, err := coll.Find(context.TODO(), filter); if err!= nil{
        return nil,err    
    }else{
        var results []*types.Notifications
        if err = cursor.All(context.TODO(), &results); err != nil {
            return nil,err
	    }else{
            return results,nil
        }
    }
}

func (s *mongoStore)SendFriendRequestNotifications(notification types.Notifications)error{
    coll := s.db.Database("eardis").Collection("notifications")
    _,err := coll.InsertOne(context.TODO(),notification); if err!=nil{
        return err
    }else{
       return nil 
    }
}

func (s *mongoStore)AcceptFriendRequest(notificationsid string, ownerid string)error{
    coll := s.db.Database("eardis").Collection("notifications")
    objectID, err := primitive.ObjectIDFromHex(notificationsid); if err != nil {return err}
    var filter bson.M
    var update bson.M
    filter = bson.M{"_id":objectID,"to":ownerid}
    var message types.Notifications
    err = coll.FindOne(context.TODO(),filter).Decode(&message); if err != nil{
        return err 
    }else{
        // add the friend to the person who sent the notification
        userColl := s.db.Database("eardis").Collection("users")
        objectID, err = primitive.ObjectIDFromHex(message.From); if err != nil {return err}
        filter = bson.M{"_id": objectID, "friends": nil}
        update = bson.M{"$set": bson.M{"friends": []string{}}}
        res , err := userColl.UpdateOne(context.TODO(), filter, update);if err!= nil{
            return err
        }else if res.MatchedCount == 0{
            filter = bson.M{"_id": objectID}
            update = bson.M{"$push": bson.M{"friends": message.To}}
	        _, err = userColl.UpdateOne(context.TODO(), filter, update);if err!= nil{
                return err
            }
        }else{
            filter = bson.M{"_id": objectID}
            update = bson.M{"$push": bson.M{"friends": message.To}}
	        _, err = userColl.UpdateOne(context.TODO(), filter, update);if err!= nil{
                return err
            }

        }
        // add the friend to the person who received the notification
        objectID, err = primitive.ObjectIDFromHex(message.To); if err != nil {return err}
        filter = bson.M{"_id": objectID, "friends": nil}
        update = bson.M{"$set": bson.M{"friends": []string{}}}
        res , err = userColl.UpdateOne(context.TODO(), filter, update);if err!= nil{
            return err
        }else if res.MatchedCount == 0{
            filter = bson.M{"_id": objectID}
            update = bson.M{"$push": bson.M{"friends": message.From}}
	        _, err = userColl.UpdateOne(context.TODO(), filter, update);if err!= nil{
                return err
            }
        }else{
            filter = bson.M{"_id": objectID}
            update = bson.M{"$push": bson.M{"friends": message.From}}
	        _, err = userColl.UpdateOne(context.TODO(), filter, update);if err!= nil{
                return err
            }

        }
        return nil
    }
}

func (s mongoStore)DeclineFriendRequest(notificationsid string,ownerid string)error{
    coll := s.db.Database("eardis").Collection("notifications")
    objectID, err := primitive.ObjectIDFromHex(notificationsid); if err != nil {return err}
    filter := bson.D{{"_id",objectID},{"to",ownerid}}
    _ , err = coll.DeleteOne(context.TODO(), filter);if err!= nil{
        return err
    }else{
        return nil
    } 
}

func (s mongoStore) GetFriendsEvents(userid string) ([]*types.Event, error){
    coll := s.db.Database("eardis").Collection("users")
    objectID, err := primitive.ObjectIDFromHex(userid); if err != nil {
		return nil,errors.New("Invalid event id") 
	}
    var filter bson.D = bson.D{{"_id",objectID}}
    var user types.User
    err = coll.FindOne(context.TODO(),filter).Decode(&user); if err != nil{
        return nil,err
    }else{
        var results []*types.Event
        eventsColl := s.db.Database("eardis").Collection("events")
        for _,item := range user.Friends{
            filter := bson.D{
		        {"owner", item},
		        {"guests", bson.D{
			        {"$in", bson.A{userid}},
		        }},
	        }

            cursor, err := eventsColl.Find(context.TODO(), filter); if err!= nil{
                return nil,err    
            }else{
                var tmp []*types.Event
                if err = cursor.All(context.TODO(), &tmp); err != nil {
                    return nil,err
	            }else{
                    results = append(results, tmp...)
                }
            }
        }
        return results,nil
    }
}

func (s mongoStore) DeleteUser(userid string)error{
    coll := s.db.Database("eardis").Collection("users")
    objectID, err := primitive.ObjectIDFromHex(userid); if err != nil {
		return errors.New("Invalid event id") 
	}
    filter := bson.D{{"_id",objectID}}
    _,err = coll.DeleteOne(context.TODO(),filter); if err != nil{
        return err
    }else{
        return nil
    }
}
