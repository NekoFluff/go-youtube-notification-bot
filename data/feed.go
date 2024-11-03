package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetFeeds() ([]ChannelFeed, error) {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("hololive-en").Collection("feeds")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var results []ChannelFeed
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetFeedByID(id primitive.ObjectID) (*ChannelFeed, error) {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("hololive-en").Collection("feeds")
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	var result ChannelFeed
	if err := collection.FindOne(context.Background(), filter).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetFeedsByName returns a list of feeds that loosely matches the name provided
func GetFeedsByName(name string, limit int) ([]ChannelFeed, error) {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("hololive-en").Collection("feeds")

	// Create a regex pattern for partial matches
	regexPattern := bson.D{primitive.E{Key: "$regex", Value: name}, primitive.E{Key: "$options", Value: "i"}}

	// Use $or to match either first name or last name
	filter := bson.D{primitive.E{Key: "$or", Value: bson.A{
		bson.D{primitive.E{Key: "firstName", Value: regexPattern}},
		bson.D{primitive.E{Key: "lastName", Value: regexPattern}},
	}}}

	// Find documents matching the filter
	cur, err := collection.Find(context.Background(), filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}

	var results []ChannelFeed
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetFeedsForUser(user string) ([]ChannelFeed, error) {
	client := GetClient()
	defer DisconnectClient(client)

	subscriptions := client.Database("hololive-en").Collection("subscriptions")

	// Aggregate pipeline to get feeds for user subscriptions
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user", Value: user}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "feeds"},
			{Key: "localField", Value: "feedID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "feedDetails"},
		}}},
		{{Key: "$unwind", Value: "$feedDetails"}},
		{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$feedDetails"}}}},
	}

	cur, err := subscriptions.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []ChannelFeed
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
