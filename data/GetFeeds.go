package data

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetFeeds() []ChannelFeed {
	client := GetClient()
	defer DisconnectClient(client)

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	// client.Database("<db>").Collection("<collection>").InsertOne(ctx, bson.D{{"x",1}})

	collection := client.Database("hololive-en").Collection("feeds")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []ChannelFeed
	if err = cur.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", results)
	return results

}
