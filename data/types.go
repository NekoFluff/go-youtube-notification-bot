package data

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ChannelFeed struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FirstName  string
	LastName   string
	TopicURL   string
	Group      string
	Generation int
}

func (cf *ChannelFeed) FullName() string {
	caser := cases.Title(language.Und)
	return strings.TrimSpace(caser.String(cf.FirstName + " " + cf.LastName))
}

type Livestream struct {
	Author  string    `bson:"author"`
	Url     string    `bson:"url"`
	Date    time.Time `bson:"date"`
	Title   string    `bson:"title"`
	Updated time.Time `bson:"updated"`
}

type Subscription struct {
	User   string             `bson:"user"`
	FeedID primitive.ObjectID `bson:"feedID"`
}

type SubscriptionGroup struct {
	User  string `bson:"_id"`
	Count int    `bson:"count"`
}
