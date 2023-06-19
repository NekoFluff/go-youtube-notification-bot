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

func GetSubscriptions(authors []string) ([]SubscriptionGroup, error) {
	client := GetClient()
	defer DisconnectClient(client)

	subscriptions := client.Database("hololive-en").Collection("subscriptions")

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "subscription", Value: bson.D{primitive.E{Key: "$in", Value: authors}}}}}}
	groupStage := bson.D{primitive.E{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$user"}, primitive.E{Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

	cur, err := subscriptions.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}

	var results []SubscriptionGroup
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetSubscriptionsForUser(user string) ([]Subscription, error) {
	client := GetClient()
	defer DisconnectClient(client)

	subscriptions := client.Database("hololive-en").Collection("subscriptions")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	filter := bson.D{primitive.E{Key: "user", Value: user}}

	cur, err := subscriptions.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []Subscription
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

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

func SaveSubscription(subscription Subscription) *mongo.UpdateResult {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("hololive-en").Collection("subscriptions")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	options := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "user", Value: subscription.User}, primitive.E{Key: "subscription", Value: subscription.Subscription}}
	update := bson.D{primitive.E{Key: "$set", Value: subscription}}

	result, err := collection.UpdateOne(ctx, filter, update, options)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Update Result: %#v\n", result)
	return result
}
