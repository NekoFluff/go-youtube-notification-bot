package pubsubhub

import (
	"encoding/xml"
)

func ParseXML(xmlString string) (Feed, error) {
	var feed Feed
	b := []byte(xmlString)

	err := xml.Unmarshal(b, &feed)

	return feed, err
}
