package main

import (
	"twint-lite/Core"
	"twint-lite/InputParser"
)

func main() {
	Arguments := InputParser.InputParser()
	Core.Main(&(Arguments.Query), &(Arguments.Instance), &(Arguments.Format))
}
