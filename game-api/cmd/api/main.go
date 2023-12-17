package main

import (
	"context"
	"fmt"
	"game-api/model"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
	mongoUrl = "mongodb://mongo:27017"
)

var client *mongo.Client

type AppConfig struct {
	mongoClient *mongo.Client
}

func main() {
	log.Println("Starting game-api service")

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconect.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	app := AppConfig{
		mongoClient: model.InitClient(mongoClient),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	//create connection options
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting: ", err)
		return nil, err
	}
	log.Println("connected to Mongo")
	return c, nil
}