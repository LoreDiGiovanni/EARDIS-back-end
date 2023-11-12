package main

import (
	"context"
	"log"
	"os"
   "errors"


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

func (s *mongoStore) getEvents(id string) ([]*Event, error){
    coll := s.db.Database("eardis").Collection("events")
    filter := bson.D{{"owner", id}}
    cursor, err := coll.Find(context.TODO(), filter); if err!= nil{
        log.Println(err,"Function: getEvents")
        return nil,err    
    }else{
        var results []*Event
        if err = cursor.All(context.TODO(), &results); err != nil {
            log.Println(err,"Function: getEvents")
            return nil,err
	    }else{
            return results,nil
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

