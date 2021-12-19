package pubsubhub

import (
	"testing"
)

func TestParseXML(t *testing.T) {
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
}
