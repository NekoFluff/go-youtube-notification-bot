package pubsubhub

import (
	"os"
	"testing"
)

func TestParseXML(t *testing.T) {
	os.Setenv("TZ", "UTC+1")

	body := `<feed xmlns:yt="http://www.youtube.com/xml/schemas/2015" xmlns="http://www.w3.org/2005/Atom">
  <link rel="hub" href="https://pubsubhubbub.appspot.com" />
  <link rel="alex" href="wtf" />
  <link rel="self" href="https://www.youtube.com/xml/feeds/videos.xml?channel_id=CHANNEL_ID" />
  <title>YouTube video feed</title>
  <updated>2015-04-01T19:05:24.552394234+00:00</updated>
  <entry>
    <id>yt:video:VIDEO_ID</id>
    <yt:videoId>VIDEO_ID</yt:videoId>
    <yt:channelId>CHANNEL_ID</yt:channelId>
    <title>Video title</title>
    <link rel="alternate" href="http://www.youtube.com/watch?v=VIDEO_ID" />
    <author>
      <name>Channel title</name>
      <uri>http://www.youtube.com/channel/CHANNEL_ID</uri>
    </author>
    <published>2015-03-06T21:40:57+00:00</published>
    <updated>2015-03-09T19:05:24.552394234+00:00</updated>
  </entry>
</feed>`
	feed, xmlError := ParseXML(body)
	if xmlError != nil {
		t.Error(xmlError)
	}

	if feed.Title != "YouTube video feed" {
		t.Errorf("Title is not the same")
	}
	entry := feed.Entries[0]
	if entry.VideoId != "VIDEO_ID" {
		t.Errorf("Video Id is incorrect")
	}
	if entry.ChannelId != "CHANNEL_ID" {
		t.Errorf("Channel Id is incorrect")
	}
	if entry.Link.Href != "http://www.youtube.com/watch?v=VIDEO_ID" {
		t.Errorf("Link is incorrect. Currently: %s", entry.Link)
	}

	currentPublishedDt := entry.Published.Format("2006-01-02 15:04:05 -0700")
	expectedPublishedDT := "2015-03-06 21:40:57 +0000"
	if currentPublishedDt != expectedPublishedDT {
		t.Errorf("Published date is incorrect. Currently: %s. Expected %s.", currentPublishedDt, expectedPublishedDT)
	}

	currentUpdatedDT := entry.Updated.Format("2006-01-02 15:04:05 -0700")
	expectedUpdatedDT := "2015-03-09 19:05:24 +0000"
	if currentUpdatedDT != expectedUpdatedDT {
		t.Errorf("Updated date is incorrect. Currently: %s. Expected: %s.", currentUpdatedDT, expectedUpdatedDT)
	}
}
