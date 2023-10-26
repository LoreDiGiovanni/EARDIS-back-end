package main

import (
	"context"
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
    log.Println("[V] New Event: ",e)
    return nil
}

func (s *mongoStore) getEvents(id string) ([]*Event, error){
    log.Println("[V] Get ",id,"'s Events")
    return nil, nil
}
