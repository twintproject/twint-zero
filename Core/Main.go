package Core

import (
)

var(
	condition bool   = true
	cursor    string = ""
)

func Main(Username *string, Instance *string) () {
	for condition {
		condition = Scrape(Request(Username, Instance, &cursor), &cursor)
	}
}