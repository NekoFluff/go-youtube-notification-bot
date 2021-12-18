package pubsubhub

import (
	"encoding/xml"
	"fmt"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Links   []Link   `xml:"link"`
	Title   string   `xml:"title"`
	Updated string   `xml:"updated"`
	Entries []Entry  `xml:"entry"`
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
}

type Entry struct {
	XMLName   xml.Name `xml:"entry"`
	Id        string   `xml:"id"`
	VideoId   string   `xml:"yt:videoId"`
	ChannelId string   `xml:"yt:channelId"`
	Title     string   `xml:"title"`
	Link      string   `xml:"link"`
	Author    Author   `xml:"author"`
	Published string   `xml:"published"`
	Updated   string   `xml:"updated"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	Uri     string   `xml:"uri"`
}

func ParseXML(xmlString string) Feed {
	var feed Feed
	b := []byte(xmlString)

	xml.Unmarshal(b, &feed)
	fmt.Println("Parsed")
	fmt.Println(feed)

	return feed
}
