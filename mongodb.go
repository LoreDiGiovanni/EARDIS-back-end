package main

import (
	"context"
	"log"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *mongoStore) initStorage() error {
    coll := s.db.Database("eardis").Collection("users")
    var result bson.M
    err := coll.FindOne(context.Background(),bson.M{}).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            log.Println("Nessun documento trovato")
        } else {
            log.Fatal("Errore nella query:", err)
            return err
        }
    }
    log.Println(result)
    return nil
}


