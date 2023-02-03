package Core

import (
	"encoding/json"
	"fmt"
)

func FormatDefault(tweet Tweet) {
	fmt.Printf("%d,%s,%s,%s,%s,%s\n",
		tweet.ID,
		tweet.URL,
		tweet.Timestamp,
		tweet.Username,
		tweet.Fullname,
		tweet.Text,
	)
}

func FormatJSON(tweet Tweet) {
	tweetJSON, _ := json.Marshal(tweet)
	fmt.Println(string(tweetJSON))
}
