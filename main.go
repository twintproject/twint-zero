package main

import (
	"twint-zero/InputParser"
	"twint-zero/Core"
)

func main() {
	Arguments := InputParser.InputParser()
	Core.Main(&(Arguments.Query), &(Arguments.Instance))
}