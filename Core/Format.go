package Core

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"os"
	
)

func FormatTweets(format string, tweets []Tweet, file *os.File) {
	if format == "json" {
		FormatTweetsJSON(tweets)
	} else {
		FormatTweetsCSV(tweets, file)
	}
}

func FormatTweetsCSV(tweets []Tweet, file *os.File) {
	var b []byte
	buf := bytes.NewBuffer(b)
	terminal := csv.NewWriter(buf)
	w := csv.NewWriter(file)


	for _, tweet := range tweets {
		
		attachments := make([]string, len(tweet.Attachments))
		for i, att := range tweet.Attachments {
			attachments[i] = *att.URL
		}
		row := []string{
			tweet.ID,
			tweet.URL,
			tweet.Timestamp,
			tweet.Username,
			tweet.Fullname,
			tweet.Text,
			strings.Join(attachments, ","),
			fmt.Sprintf("%d", tweet.Stats.Replies),
			fmt.Sprintf("%d", tweet.Stats.Retweets),
			fmt.Sprintf("%d", tweet.Stats.Quotes),
			fmt.Sprintf("%d", tweet.Stats.Likes),
		}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing row to csv:", err)
		}

		if err := terminal.Write(row); err != nil {
			log.Fatalln("error writing row to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	terminal.Flush()
	if err := terminal.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(buf.Bytes()))

}

func FormatTweetsJSON(tweets []Tweet) {
	for _, tweet := range tweets {
		tweetJSON, _ := json.Marshal(tweet)
		fmt.Println(string(tweetJSON))
	}
}
