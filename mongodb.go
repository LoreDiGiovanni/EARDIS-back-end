package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/bson"
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
    coll := s.db.Database("eardis").Collection("users")
    _,err := coll.InsertOne(context.TODO(),user)
    if err!=nil{
        return nil,err 
    }else{
        jwt,err:= createUserJWT(user)
        if err!= nil{
            return nil,err
        }else{
            user.JWT = jwt
            return user,nil
        }
    }
}

