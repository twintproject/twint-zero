package Core

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	escapeChars *strings.Replacer = strings.NewReplacer(
		"\n", " ",
		"\t", " ",
		"\r", " ",
		";", " ",
		",", " ",
	)
	Webpage *goquery.Document = new(goquery.Document)
)

type Tweet struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Text      string `json:"text"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Timestamp string `json:"timestamp"`
}

func extractViaRegexp(text *string, re string) string {
	theRegex := regexp.MustCompile(re)
	match := theRegex.Find([]byte(*text))
	return string(match[:])
}

func Scrape(responseBody io.ReadCloser, Format *string, cursor *string) bool {
	parsedWebpage, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		log.Fatal("[x] cannot parse webpage. Please report to admins with the query attached.")
	}
	defer responseBody.Close()

	if parsedWebpage.Find("div.timeline-footer").Length() > 0 {
		return false
	}

	var tweets []Tweet
	parsedWebpage.Find("div.timeline-item").Each(func(i int, t *goquery.Selection) {
		tweet_ID_h, _ := t.Find("a").Attr("href")
		tweet_ID_s := strings.Split(tweet_ID_h, "/")
		tweet_ID, _ := strconv.Atoi(extractViaRegexp(&(tweet_ID_s[len(tweet_ID_s)-1]), `\d*`))

		tweet_URL := fmt.Sprintf("https://twitter.com%s", strings.Split(tweet_ID_h, "#")[0])

		tweet_TS, _ := t.Find("span.tweet-date").Find("a").Attr("title")

		tweet_text := escapeChars.Replace(t.Find("div.tweet-content.media-body").Text())

		tweet_handle := t.Find("a.username").First().Text()
		tweet_fname := t.Find("a.fullname").First().Text()

		if tweet_ID != 0 {
			tweet := Tweet{
				ID:        tweet_ID,
				URL:       tweet_URL,
				Text:      tweet_text,
				Username:  tweet_handle,
				Fullname:  tweet_fname,
				Timestamp: tweet_TS,
			}
			tweets = append(tweets, tweet)
		}
	})

	FormatTweets(*Format, tweets)

	*cursor, _ = parsedWebpage.Find("div.show-more").Last().Find("a").Attr("href")
	return true
}
