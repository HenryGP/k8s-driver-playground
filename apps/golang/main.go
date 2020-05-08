package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kr/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	mongoURI := os.Getenv("URI")

	if len(mongoURI) == 0 {
		panic("No URI")
	}

	poolMonitor := &event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			log.Printf("%# v", pretty.Formatter(evt))
		},
	}

	// Don't set the Command monitor if it gets too verbose for the CRUD operations
	/*cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Printf("%# v", pretty.Formatter(evt))
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			log.Printf("%# v", pretty.Formatter(evt))
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			log.Printf("%# v", pretty.Formatter(evt))
		},
	}*/

	clientOptions := options.Client().ApplyURI(mongoURI).SetPoolMonitor(poolMonitor) //.SetMonitor(cmdMonitor)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	database := client.Database("golang")
	collection := database.Collection("demo")

	log.Println("Dropping the sample collection")
	_ = database.RunCommand(
		context.Background(),
		bson.D{{"drop", "demo"}},
	)

	for {
		_, err = collection.InsertOne(ctx, bson.D{{"x", 1}})
		time.Sleep(1 * time.Second)
		cur, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(context.Background())
		for cur.Next(context.Background()) {
			doc := bson.D{}
			err := cur.Decode(&doc)
			if err != nil {
				fmt.Println("Failed to decode bytes")
				log.Fatal(err)
			}
			// Uncomment to print doc
			// fmt.Println(doc)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
		collection.DeleteOne(ctx, bson.D{{"x", 1}})
	}
}
