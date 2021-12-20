package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveLivestream(livestream Livestream) *mongo.UpdateResult {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("hololive-en").Collection("scheduledLivestreams")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	options := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "url", Value: livestream.Url}}
	update := bson.D{primitive.E{Key: "$set", Value: livestream}}

	result, err := collection.UpdateOne(ctx, filter, update, options)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Update Result: %#v", result)
	return result
}
