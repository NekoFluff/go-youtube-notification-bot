package pubsubhub

import (
	"encoding/xml"
	"fmt"
)

func ParseXML(xmlString string) Feed {
	var feed Feed
	b := []byte(xmlString)

	xml.Unmarshal(b, &feed)
	fmt.Println("Parsed")
	fmt.Println(feed)

	return feed
}
