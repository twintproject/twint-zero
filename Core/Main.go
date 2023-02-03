package Core

import (
	"net/url"
)

var (
	condition bool   = true
	cursor    string = ""
)

func Main(Query *string, Instance *string, Format *string) {
	(*Query) = url.QueryEscape(*Query)
	for condition {
		condition = Scrape(Request(Query, Instance, &cursor), Format, &cursor)
	}
}
