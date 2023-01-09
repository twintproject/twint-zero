package Core

import (
	"io"
	"fmt"
	"log"
	"time"
	"net/http"
)

var (
	Client *http.Client = new(http.Client);
)

func Request(Query *string, Instance *string, cursor *string) (io.ReadCloser) {
	var url string = fmt.Sprintf("https://%s/search?f=tweet&q=%s", *Instance, *Query)
	if *cursor != "" {
		url = fmt.Sprintf("https://%s/search%s", *Instance, *cursor)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("[nr] %s\n", err)
	}

	req.Header.Set("User-Agent", "twint-zero")
	res, err := Client.Do(req)
	if err != nil {
		log.Fatalf("[do] %s\n", err)
	}

	if res.StatusCode != 200 {
		if 500 <= res.StatusCode && res.StatusCode <= 599 {
			time.Sleep(10 * time.Second)
			return Request(Query, Instance, cursor)
		} else {
			log.Fatalf("status code error: %d %s \n %s", res.StatusCode, res.Status, url)
		}
	}
	return res.Body
}