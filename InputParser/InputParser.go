package InputParser

import (
	"flag"
	"os"
)

type Arguments struct {
	Query    string
	Instance string
	Format   string
	FileName string
}

var arguments *Arguments = new(Arguments)

func InputParser() *Arguments {

	flag.StringVar(&(arguments.Query), "Query", "", "Specify search query.")
	flag.StringVar(&(arguments.Instance), "Instance", "nitter.nl", "Specify instance to get data from.")
	flag.StringVar(&(arguments.Format), "Format", "csv", "Specify the return format: csv (default), or json.")
	flag.StringVar(&(arguments.FileName), "FileName", "data", "Specify csv output file name: data (default).")
	flag.Parse()

	if (*arguments).Query == "" || !ValidateFormatArgument(arguments) {
		flag.Usage()
		os.Exit(1)
	}

	return arguments
}

func ValidateFormatArgument(arguments *Arguments) bool {
	format := (*arguments).Format
	return format == "" || format == "csv" || format == "json"
}
