package InputParser

import (
	"os"
	"flag"
)

type Arguments struct {
	Query    string
	Instance string
}

var arguments *Arguments = new(Arguments)

func InputParser() *Arguments {

	flag.StringVar(&(arguments.Query), "Query", "", "Specify search query.")
	flag.StringVar(&(arguments.Instance), "Instance", "nitter.nl", "Specify instance to get data from.")
    flag.Parse()

	if (*arguments).Query == "" { 
		flag.Usage()
		os.Exit(1)
	}	
	
	return arguments
}