package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
