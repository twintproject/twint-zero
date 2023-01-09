package main

import (
	"twint-lite/InputParser"
	"twint-lite/Core"
)

func main() {
	Arguments := InputParser.InputParser()
	Core.Main(&(Arguments.Query), &(Arguments.Instance))
}