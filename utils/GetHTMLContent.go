package utils

import (
	"io/ioutil"
	"log"
	"net/http"
)

func GetHTMLContent(url string) (html []byte, err error) {
	// Get html content
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// Reads html as a slice of bytes
	html, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}
