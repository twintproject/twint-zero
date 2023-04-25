package Core

import (
	"log"
	"net/url"
	"os"
)

var (
	condition bool   = true
	cursor    string = ""
)

func Main(Query *string, Instance *string, Format *string, FileName *string) {
	(*Query) = url.QueryEscape(*Query)
	name := *FileName
	name += ".csv"
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	for condition {
		condition = Scrape(Request(Query, Instance, &cursor), Format, &cursor, file)
	}
}
