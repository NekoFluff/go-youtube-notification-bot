package pubsubhub

import (
	"fmt"
	"testing"
	"time"
)

func TestGetLivestreamUnixTime(t *testing.T) {
	unixTime, err := GetLivestreamUnixTime("https://www.youtube.com/watch?v=c7K6RInG3Dw")

	if err != nil {
		t.Error(err)
	}

	// TODO: Maybe remove this check and just rely on the above error not being present
	str := fmt.Sprintf("%v", unixTime)
	if str != "2022-09-28 13:00:00 -0700 MST" {
		t.Errorf("Unix Timestamp: %s", str)
	}
}

func TestConvertEntryToLivestream(t *testing.T) {
	entry := Entry{
		Id:        "id",
		VideoId:   "video id",
		ChannelId: "channel id",
		Title:     "title",
		Link: Link{
			Href: "https://www.youtube.com/watch?v=c7K6RInG3Dw",
		},
		Author: Author{
			Name: "author name",
			Uri:  "uri",
		},
		Published: time.Now(),
		Updated:   time.Now(),
	}

	livestream, err := ConvertEntryToLivestream(entry)
	if err != nil {
		t.Error(err)
	}

	if livestream.Author != entry.Author.Name {
		t.Errorf("Author is not correct. Currently %s", livestream.Author)
	}
	if livestream.Url != entry.Link.Href {
		t.Error("Url is not correct")
	}
	if livestream.Title != entry.Title {
		t.Error("Title is not correct")
	}
	if livestream.Updated != entry.Updated {
		t.Error("Title is not correct")
	}
}
