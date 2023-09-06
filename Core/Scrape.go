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

type Tweet struct {
	ID          string       `json:"id"`
	URL         string       `json:"url"`
	Text        string       `json:"text"`
	ReplyTo 	string 		 `json:"reply_to,omitempty"`
	Username    string       `json:"username"`
	Fullname    string       `json:"fullname"`
	Timestamp   string       `json:"timestamp"`
	Attachments []Attachment `json:"attachments"`

	Stats TweetStats `json:"stats"`

	QuoteFullname string `json:"quote_fullname,omitempty"`
    	QuoteUsername string `json:"quote_username,omitempty"`
    	QuoteDate     string `json:"quote_date,omitempty"`
    	QuoteID       string `json:"quote_id,omitempty"`
    	QuoteText     string `json:"quote_text,omitempty"`
}

type Attachment struct {
	Type            string  `json:"type"`
	URL             *string `json:"url"`
	PreviewImageURL *string `json:"preview_image_url"`
	AltText         *string `json:"alt_text"`
}

type TweetStats struct {
	Replies  int64 `json:"replies"`
	Retweets int64 `json:"retweets"`
	Quotes   int64 `json:"quotes"`
	Likes    int64 `json:"likes"`
}

func extractViaRegexp(text *string, re string) string {
	theRegex := regexp.MustCompile(re)
	match := theRegex.Find([]byte(*text))
	return string(match[:])
}

func Scrape(responseBody io.ReadCloser, Instance *string, Format *string, cursor *string) bool {
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
		tweet_ID := extractViaRegexp(&(tweet_ID_s[len(tweet_ID_s)-1]), `\d*`)

		tweet_URL := fmt.Sprintf("https://twitter.com%s", strings.Split(tweet_ID_h, "#")[0])

		tweet_TS, _ := t.Find("span.tweet-date").Find("a").Attr("title")

		tweet_text := t.Find("div.tweet-content.media-body").Text()

		reply_to := t.Find("div.replying-to").Text()

		tweet_handle := t.Find("a.username").First().Text()
		tweet_fname := t.Find("a.fullname").First().Text()

		tweet_stats := t.Find("div.tweet-stats")
		tweet_stats_reply, _ := strconv.ParseInt(
			strings.TrimSpace(
				strings.ReplaceAll(
					tweet_stats.Find("span.tweet-stat").Eq(0).Text(), ",", "",)), 10, 64)
		tweet_stats_retweet, _ := strconv.ParseInt(
			strings.TrimSpace(
				strings.ReplaceAll(
					tweet_stats.Find("span.tweet-stat").Eq(1).Text(), ",", "")), 10, 64)
		tweet_stats_quote, _ := strconv.ParseInt(
			strings.TrimSpace(
				strings.ReplaceAll(
					tweet_stats.Find("span.tweet-stat").Eq(2).Text(), ",", "")), 10, 64)
		tweet_stats_like, _ := strconv.ParseInt(
			strings.TrimSpace(
				strings.ReplaceAll(
					tweet_stats.Find("span.tweet-stat").Eq(3).Text(), ",", "")), 10, 64)
		
		tweet_attachments := make([]Attachment, 0)
		t.Find("div.attachments").Find("div.attachment.image").Find("img").Each(func(i int, s *goquery.Selection) {
			src, exists := s.Attr("src")
			alt, _ := s.Attr("alt")
			if exists {
				src = fmt.Sprintf("https://%s%s", *Instance, src)
				tweet_attachments = append(tweet_attachments, Attachment{
					Type:    "photo",
					URL:     &src,
					AltText: &alt,
				})
			}
		})
		t.Find("div.attachments").Find("video.gif").Each(func(i int, s *goquery.Selection) {
			preview, exists := s.Attr("poster")
			if exists {
				src, _ := s.Find("source").Attr("src")
				preview = fmt.Sprintf("https://%s%s", *Instance, preview)
				src = fmt.Sprintf("https://%s%s", *Instance, src)
				tweet_attachments = append(tweet_attachments, Attachment{
					Type:            "animated_gif",
					URL:             &src,
					PreviewImageURL: &preview,
				})
			}
		})
		t.Find("div.attachments").Find("div.gallery-video").Find("video").Each(func(i int, s *goquery.Selection) {
			preview, exists := s.Attr("poster")
			if exists {
				var ur *string
				src, exists := s.Attr("data-url")
				src = fmt.Sprintf("https://%s%s", *Instance, src)
				preview = fmt.Sprintf("https://%s%s", *Instance, preview)
				if exists {
					ur = &src
				}
				tweet_attachments = append(tweet_attachments, Attachment{
					Type:            "video",
					URL:             ur,
					PreviewImageURL: &preview,
				})
			}
		})

		var quote_fullname, quote_username, quote_date, quote_text, quote_id string
		if quoteInfo := t.Find("div.quote"); quoteInfo.Length() > 0 {
			quote_fullname = quoteInfo.Find("a.fullname").Text()
			quote_username = quoteInfo.Find("a.username").Text()
			quote_date, _ = quoteInfo.Find("span.tweet-date").Find("a").Attr("title")
			quote_id_h, _ := quoteInfo.Find("a.quote-link").Attr("href")
			quote_id_s := strings.Split(quote_id_h, "/")
			quote_id = extractViaRegexp(&(quote_id_s[len(quote_id_s)-1]), `\d*`)
			quote_text = quoteInfo.Find("div.quote-text").Text()
		}

		stats := TweetStats{
			Replies:  tweet_stats_reply,
			Retweets: tweet_stats_retweet,
			Quotes:   tweet_stats_quote,
			Likes:    tweet_stats_like,
		}

		if tweet_ID != "" {
			tweet := Tweet{
				ID:          tweet_ID,
				URL:         tweet_URL,
				Text:        tweet_text,
				ReplyTo: 	 reply_to,
				Username:    tweet_handle,
				Fullname:    tweet_fname,
				Timestamp:   tweet_TS,
				Attachments: tweet_attachments,
				Stats:       stats,
				QuoteFullname: quote_fullname,
		                QuoteUsername: quote_username,
		                QuoteDate:     quote_date, 
		                QuoteID:       quote_id,
		                QuoteText:     quote_text,
			}
			tweets = append(tweets, tweet)
		}
	})

	if len(tweets) == 0 {
		return false
	}

	FormatTweets(*Format, tweets)

	*cursor, _ = parsedWebpage.Find("div.show-more").Last().Find("a").Attr("href")
	return true
}
