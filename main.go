package main

import (
	"twint-zero/Core"
	"twint-zero/InputParser"
)

func main() {
	Arguments := InputParser.InputParser()
	Core.Main(&(Arguments.Query), &(Arguments.Instance), &(Arguments.Format))
}
