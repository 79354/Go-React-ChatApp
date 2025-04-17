package config

import (
	"context"
	"os"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDatabase(){
	log.Println("Connecting to database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbURI := os.Getenv("DB_URL")
	options := options.Client().ApplyURI(dbURI)
	if dbURI == ""{
		dbURI = "27017"
	}

	client, err := mongo.Connect(ctx, options)
	if err != nil{
		log.Fatal("Error connecting to mongodb", err)
	}

	if client.Ping(ctx, nil); err != nil{
		log.Fatal("Can't ping the client", err)
	}

	Client = client
	log.Println("Connected to database")
}

func DisConnectDB(){
	if Client != nil{
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := Client.Disconnect(ctx); err != nil{
			log.Println("Error disconnecting MongoDB:", err)
		}else {
			log.Println("Database disconnected successfully.")
		}
	}
}