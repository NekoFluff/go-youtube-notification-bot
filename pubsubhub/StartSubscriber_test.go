package pubsubhub

import (
	"fmt"
	"testing"
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
		Published: "published",
		Updated:   "updated",
	}

	livestream, err := ConvertEntryToLivestream(entry)
	if err != nil {
		t.Error(err)
	}

	if livestream.Author != "author name" {
		t.Errorf("Author is not correct. Currently %s", livestream.Author)
	}
	if livestream.Url != "https://www.youtube.com/watch?v=c7K6RInG3Dw" {
		t.Error("Url is not correct")
	}
	if livestream.Title != "title" {
		t.Error("Title is not correct")
	}
}
