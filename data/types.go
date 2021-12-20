package data

import (
	"time"
)

type ChannelFeed struct {
	FirstName  string
	LastName   string
	TopicURL   string
	Group      string
	Generation int
}

type Livestream struct {
	Author  string    `bson:"author"`
	Url     string    `bson:"url"`
	Date    time.Time `bson:"date"`
	Title   string    `bson:"title"`
	Updated time.Time `bson:"updated"`
}
