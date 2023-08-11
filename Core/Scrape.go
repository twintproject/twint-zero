package Core

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Tweet struct {
	ID          string       `json:"id"`
	URL         string       `json:"url"`
	Text        string       `json:"text"`
	Username    string       `json:"username"`
	Fullname    string       `json:"fullname"`
	Timestamp   string       `json:"timestamp"`
	Attachments []Attachment `json:"attachments"`
        ProfileImageURL string `json:"profile_image_url"`
	Stats TweetStats `json:"stats"`
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

		tweet_handle := t.Find("a.username").First().Text()
		tweet_fname := t.Find("a.fullname").First().Text()

                // Fetch the profile image URL using ScrapeUserProfile function
                profileImageURL := ScrapeUserProfile(tweet_handle, *Instance)
		// tweet stats: reply, retweet, quote, like as span.tweet-stats childs of tweet_stats
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
				Username:    tweet_handle,
				Fullname:    tweet_fname,
				Timestamp:   tweet_TS,
				Attachments: tweet_attachments,
				Stats:       stats,
                                ProfileImageURL: profileImageURL,  // Add this line
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
// ScrapeUserProfile fetches the user's profile on Nitter and extracts the profile image URL.
func ScrapeUserProfile(username string, instance string) string {
        username = strings.TrimPrefix(username, "@")


	// Construct the URL to the user's profile on the instance (e.g., Nitter)
	profileURL := fmt.Sprintf("https://%s/%s", instance, username)

	// Fetch the profile page using http.Get
	res, err := http.Get(profileURL)
	if err != nil {
		log.Printf("Error fetching profile for user %s: %v", username, err)
		return ""
	}
	defer res.Body.Close()

	// Parse the response body with goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Error parsing profile for user %s: %v", username, err)
		return ""
	}

	// Extract the profile image URL
	relativeProfileImageURL, exists := doc.Find("a.profile-card-avatar").Attr("href")
	if !exists {
		log.Printf("Profile image not found for user %s", username)
		return ""
	}

	// Construct the full profile image URL
	return fmt.Sprintf("https://%s%s", instance, relativeProfileImageURL)
}

