package main

import (
	"context"
	"errors"
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

func newMongoStore() (*mongoStore,error,){
    uri := os.Getenv("MONGODB_URI")
    client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
    if err!= nil {
        return nil,err
    }else{
        return &mongoStore{db: client},nil
    }
}

func (s *mongoStore) createAccount(user *User) (*User,error){
    var err error
    coll := s.db.Database("eardis").Collection("users")
    res,err := coll.InsertOne(context.TODO(),user)
    if err != nil{
        return nil,err
    }else{
        id := res.InsertedID.(primitive.ObjectID)
        user.ID = id.Hex()
        user.JWT,err= createUserJWT(user)
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

func (s *mongoStore) createEvent(e *Event) error{
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

func (s *mongoStore) getEvents(userid string) ([]*Event, error){
    coll := s.db.Database("eardis").Collection("events")
    filter := bson.D{{"owner", userid}}
    cursor, err := coll.Find(context.TODO(), filter); if err!= nil{
        return nil,err    
    }else{
        var results []*Event
        if err = cursor.All(context.TODO(), &results); err != nil {
            return nil,err
	    }else{
            fireds_resout,err := s.getFriendsEvents(userid);if err != nil{
                return results,nil
            }else{
                results := append(results, fireds_resout...)
                return results,nil
            }
        }
    }
}
func (s *mongoStore) deleteEvent(ownerid string ,eventid string) error{
    coll := s.db.Database("eardis").Collection("events")
    objectID, err := primitive.ObjectIDFromHex(eventid); if err != nil {return err}
    filter := bson.D{{"_id", objectID},{"owner", ownerid}}
    _, err = coll.DeleteOne(context.TODO(), filter); if err != nil {
        return err
    }else{
        return nil
    }
}

func (s *mongoStore) patchEvent(ownerid string ,eventid string,e *Event) error{
    coll := s.db.Database("eardis").Collection("events")
    objectID, err := primitive.ObjectIDFromHex(eventid); if err != nil {
		return errors.New("Invalid event id") 
	}
    filter := bson.D{{"_id", objectID},{"owner", ownerid}}
    update := bson.D{{"$set", e}}
    _ , err = coll.UpdateOne(context.TODO(), filter, update); if err != nil {
        return err
    }else{
        return nil
    }
}

func (s *mongoStore)login(user *User) (string,error){
    coll := s.db.Database("eardis").Collection("users")
    filter := bson.D{{"username",user.Username},{"pwd",user.PWD},{"email",user.Email}}
    var newuser User
    err := coll.FindOne(context.TODO(),filter).Decode(&newuser); if err != nil{
        return "", err 
    }else{
        return newuser.JWT, nil
    }
}

func (s *mongoStore)searchUser(email string) (*DisplayableUser,error){
    coll := s.db.Database("eardis").Collection("users")
    filter := bson.D{{"email",email}}
    var friend User
    err := coll.FindOne(context.TODO(),filter).Decode(&friend); if err != nil{
        return nil,err 
    }else{
        var duser DisplayableUser = DisplayableUser{ID: friend.ID,Username: friend.Username,Email: friend.Email}
        return &duser, nil
    }
}

func (s *mongoStore)getUser(userid string)(*DisplayableUser,error){
    coll := s.db.Database("eardis").Collection("users")
    objectID, err := primitive.ObjectIDFromHex(userid); if err != nil {
		return nil,errors.New("Invalid event id") 
	}
    filter := bson.D{{"_id",objectID}}
    var user User
    err = coll.FindOne(context.TODO(),filter).Decode(&user); if err != nil{
        return nil,err 
    }else{
        var duser DisplayableUser = DisplayableUser{ID: user.ID,Username: user.Username,Email: user.Email}
        return &duser, nil
    }
}

func (s *mongoStore)getNotifications(userid string)([]*Notifications,error){    
    coll := s.db.Database("eardis").Collection("notifications")
    filter := bson.D{{"to", userid}}
    cursor, err := coll.Find(context.TODO(), filter); if err!= nil{
        return nil,err    
    }else{
        var results []*Notifications
        if err = cursor.All(context.TODO(), &results); err != nil {
            return nil,err
	    }else{
            return results,nil
        }
    }
}

func (s *mongoStore)sendFriendRequestNotifications(notification Notifications)error{
    coll := s.db.Database("eardis").Collection("notifications")
    _,err := coll.InsertOne(context.TODO(),notification); if err!=nil{
        return err
    }else{
       return nil 
    }
}

func (s *mongoStore)acceptFriendRequest(notificationsid string, ownerid string)error{
    coll := s.db.Database("eardis").Collection("notifications")
    objectID, err := primitive.ObjectIDFromHex(notificationsid); if err != nil {return err}
    var filter bson.M
    var update bson.M
    filter = bson.M{"_id":objectID,"to":ownerid}
    var message Notifications
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

func (s mongoStore)declineFriendRequest(notificationsid string,ownerid string)error{
    coll := s.db.Database("eardis").Collection("notifications")
    objectID, err := primitive.ObjectIDFromHex(notificationsid); if err != nil {return err}
    filter := bson.D{{"_id",objectID},{"to",ownerid}}
    _ , err = coll.DeleteOne(context.TODO(), filter);if err!= nil{
        return err
    }else{
        return nil
    } 
}

func (s mongoStore) getFriendsEvents(userid string) ([]*Event, error){
    coll := s.db.Database("eardis").Collection("users")
    objectID, err := primitive.ObjectIDFromHex(userid); if err != nil {
		return nil,errors.New("Invalid event id") 
	}
    var filter bson.D = bson.D{{"_id",objectID}}
    var user User
    err = coll.FindOne(context.TODO(),filter).Decode(&user); if err != nil{
        return nil,err
    }else{
        var results []*Event
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
                var tmp []*Event
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
