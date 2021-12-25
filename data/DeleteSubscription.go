package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteSubscription(subscription Subscription) *mongo.DeleteResult {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("hololive-en").Collection("subscriptions")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	filter := bson.D{primitive.E{Key: "user", Value: subscription.User}, primitive.E{Key: "subscription", Value: subscription.Subscription}}

	result, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Delete Result: %#v\n", result)
	return result
}
