package data

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLivestream(url string) (*Livestream, error) {
	client := GetClient()
	defer DisconnectClient(client)

	subscriptions := client.Database("hololive-en").Collection("scheduledLivestreams")
	var livestream Livestream
	filter := bson.D{primitive.E{Key: "url", Value: url}}
	if err := subscriptions.FindOne(context.Background(), filter).Decode(&livestream); err != nil {
		return nil, err
	}

	return &livestream, nil
}

func GetLivestreams() ([]Livestream, error) {
	client := GetClient()
	defer DisconnectClient(client)

	subscriptions := client.Database("hololive-en").Collection("scheduledLivestreams")
	cur, err := subscriptions.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var results []Livestream
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

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
		slog.Error("Failed to save livestream", "error", err)
		return nil
	}

	fmt.Printf("Update Result: %#v\n", result)
	return result
}
