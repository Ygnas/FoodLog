package controllers

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

type FirebaseDatabase struct {
	*db.Client
	Storage *storage.Client
}

var firebaseDatabase FirebaseDatabase

func (db *FirebaseDatabase) FirebaseConnect() error {
	ctx := context.Background()

	// databaseURL := os.Getenv("DATABASE_URL")

	opt := option.WithCredentialsFile("foodlog-credentials.json")
	conf := &firebase.Config{
		DatabaseURL:   "https://foodlog-9c3fd-default-rtdb.europe-west1.firebasedatabase.app/",
		StorageBucket: "foodlog-9c3fd.appspot.com",
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firebase database: %v\n", err)
	}

	storeClient, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firebase Storage: %v\n", err)
	}

	db.Client = client
	db.Storage = storeClient
	return nil
}

func GetFirebaseDatabase() *FirebaseDatabase {
	return &firebaseDatabase
}
