package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

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
