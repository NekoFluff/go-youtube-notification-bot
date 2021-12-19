package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetLivestreams() []Livestream {
	client := GetClient()
	defer DisconnectClient(client)

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	subscriptions := client.Database("hololive-en").Collection("scheduledLivestreams")
	cur, err := subscriptions.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []Livestream
	if err = cur.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	return results
}
