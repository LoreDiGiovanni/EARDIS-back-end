package main

import (
    "context"
    "os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore struct{
    db *mongo.Client
}

func newMongoStore() (error,*mongo.Client){
    uri := os.Getenv("MONGODB_URI")
    client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
    if err!= nil {
        return err,nil
    }else{
        return nil,client
    }
}
