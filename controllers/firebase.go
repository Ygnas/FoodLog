package controllers

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

type FirebaseDatabase struct {
	*db.Client
}

var firebaseDatabase FirebaseDatabase

func (db *FirebaseDatabase) FirebaseConnect() error {
	ctx := context.Background()

	opt := option.WithCredentialsFile("foodlog-credentials.json")
	conf := &firebase.Config{
		DatabaseURL: "https://foodlog-9c3fd-default-rtdb.europe-west1.firebasedatabase.app/",
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firebase database: %v\n", err)
	}

	db.Client = client
	return nil
}

func GetFirebaseDatabase() *FirebaseDatabase {
	return &firebaseDatabase
}
