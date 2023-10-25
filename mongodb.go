package main

import (
	"context"
	"log"
	"os"

	//"go.mongodb.org/mongo-driver/bson"
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
    user.JWT,err= createUserJWT(user)
    res,err := coll.InsertOne(context.TODO(),user)
    if err != nil{
        return nil,err
    }else{
        id := res.InsertedID.(primitive.ObjectID)
        user.ID = string(id.String())
        return user,nil
    }
}


func (s *mongoStore) getEvent(*User) (*Event, error){
    return nil, nil
}
