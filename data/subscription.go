package data

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSubscriptions(authors []string) ([]SubscriptionGroup, error) {
	client := GetClient()
	defer DisconnectClient(client)

	feeds := client.Database("hololive-en").Collection("feeds")

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "firstName", Value: bson.D{{Key: "$in", Value: authors}}}},
		bson.D{{Key: "lastName", Value: bson.D{{Key: "$in", Value: authors}}}},
	}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "subscriptions"}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "feedID"}, {Key: "as", Value: "subscriptions"}}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$subscriptions"}}}}
	replaceStage := bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$subscriptions"}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$user"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}

	cur, err := feeds.Aggregate(context.Background(), mongo.Pipeline{matchStage, lookupStage, unwindStage, replaceStage, groupStage})
	if err != nil {
		return nil, err
	}

	var results []SubscriptionGroup
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("Found %d subscriptions to notify for the authors %v\n", len(results), authors))

	return results, nil
}

func GetSubscriptionsForUser(user string) ([]Subscription, error) {
	client := GetClient()
	defer DisconnectClient(client)

	subscriptions := client.Database("hololive-en").Collection("subscriptions")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	filter := bson.D{{Key: "user", Value: user}}

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
	filter := bson.D{{Key: "user", Value: subscription.User}, {Key: "feedID", Value: subscription.FeedID}}

	result, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		slog.Error("Failed to delete subscription", "error", err)
		return nil
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
	filter := bson.D{{Key: "user", Value: subscription.User}, {Key: "feedID", Value: subscription.FeedID}}
	update := bson.D{{Key: "$set", Value: subscription}}

	result, err := collection.UpdateOne(ctx, filter, update, options)

	if err != nil {
		slog.Error("Failed to save subscription", "error", err)
	}

	fmt.Printf("Update Result: %#v\n", result)
	return result
}
