package pubsubhub

import (
	"testing"
)

func TestOldGetLivestreamUnixTime(t *testing.T) {
	_, err := OldGetLivestreamUnixTime("https://www.youtube.com/watch?v=c7K6RInG3Dw")

	if err != nil {
		t.Error(err)
	}
}

// func TestConvertEntryToLivestream(t *testing.T) {
// 	entry := Entry{
// 		Id:        "id",
// 		VideoId:   "video id",
// 		ChannelId: "channel id",
// 		Title:     "title",
// 		Link: Link{
// 			Href: "https://www.youtube.com/watch?v=c7K6RInG3Dw",
// 		},
// 		Author: Author{
// 			Name: "author name",
// 			Uri:  "uri",
// 		},
// 		Published: time.Now(),
// 		Updated:   time.Now(),
// 	}

// 	livestream, err := ConvertEntryToLivestream(entry)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if livestream.Author != entry.Author.Name {
// 		t.Errorf("Author is not correct. Currently %s", livestream.Author)
// 	}
// 	if livestream.Url != entry.Link.Href {
// 		t.Error("Url is not correct")
// 	}
// 	if livestream.Title != entry.Title {
// 		t.Error("Title is not correct")
// 	}
// 	if livestream.Updated != entry.Updated {
// 		t.Error("Title is not correct")
// 	}
// }
