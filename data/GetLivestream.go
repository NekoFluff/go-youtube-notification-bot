package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
