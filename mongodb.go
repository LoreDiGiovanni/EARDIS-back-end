package main

import (
	"context"
	"fmt"
	"os"

	//"go.mongodb.org/mongo-driver/bson"
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
    coll := s.db.Database("eardis").Collection("users")
    var result User
    err := coll.FindOne(context.TODO(),user).Decode(&result)
    if err == mongo.ErrNoDocuments{
        jwt,err:= createUserJWT(user)
        if err != nil {
            return nil,fmt.Errorf("Impossible to create token")
        }else{
            user.JWT = jwt
            _,err := coll.InsertOne(context.TODO(),user)
            if err != nil{
                return nil,fmt.Errorf("Impossible to create account")
            }else{
                coll.FindOne(context.TODO(),user).Decode(&result)
                user.ID = result.ID;
                return user,nil
            }
        }
	}else{
        return nil,fmt.Errorf("User exist!")
    }
}

func (s *mongoStore) getEvent(*User) (*Event, error){
    return nil, nil
}
