package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MyAggResult struct {
	Id         string   `bson:"_id"`
	TotalTips  float64  `bson:"total_tips"`
	PickupArea []string `bson:"pickup_area"`
}

func randomInsert()

func main() {
	mongoURI := os.Getenv("URI")

	if len(mongoURI) == 0 {
		panic("No URI")
	}

	iteration := 0

	for {

		client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Iteration %v\n", iteration)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			panic(err)
		}

		database := client.Database("golang")
		collection := database.Collection("tour")

		log.Println("Dropping collection golang.tour (command)")
		result := database.RunCommand(
			context.Background(),
			bson.D{{"drop", "tour"}},
		)
		fmt.Println(result)

		fmt.Println("Inserting a single document.")
		res1, err := collection.InsertOne(
			context.Background(),
			bson.D{
				{"name", "MongoDB"},
				{"type", "database"},
				{"count", 1},
				{"tags", bson.A{"webscale"}},
				{"info", bson.D{
					{"x", 203},
					{"y", 102},
					{"z", "N/A"},
				}},
			})

		if err != nil {
			log.Fatal(err)
		}
		id := res1.InsertedID
		fmt.Println(id)

		fmt.Println("Inserting multiple documents.")
		_, err = collection.InsertMany(
			context.Background(),
			[]interface{}{
				bson.D{
					{"item", "journal"},
					{"qty", 25},
					{"tags", bson.A{"blank", "red"}},
					{"size", bson.D{
						{"h", 14},
						{"w", 21},
						{"uom", "cm"},
					}},
				},
				bson.D{
					{"item", "mat"},
					{"qty", 25},
					{"tags", bson.A{"gray"}},
					{"size", bson.D{
						{"h", 27.9},
						{"w", 35.5},
						{"uom", "cm"},
					}},
				},
			})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Iterating over a collection with Find()")
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
			fmt.Println(doc)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Running Aggregation Pipeline on taxidata.chicago")
		collection = client.Database("taxidata").Collection("chicago")

		pipeline := mongo.Pipeline{
			{{"$match", bson.D{
				{"trip_seconds", bson.D{{"$gt", 1000}}},
			}}},
			{{"$group", bson.D{
				{"_id", "$payment_type"},
				{"total_tips", bson.D{{"$sum", "$tips"}}},
				{"pickup_area", bson.D{{"$addToSet", "$pickup_community_area"}}},
			}}},
		}
		fmt.Println(pipeline)
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			doc := new(MyAggResult)
			err := cursor.Decode(doc)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(doc)
		}
		fmt.Println("Done.")

		fmt.Println("Sleeping for 2 seconds")
		time.Sleep(2 * time.Second)
		iteration += 1
	}
}
