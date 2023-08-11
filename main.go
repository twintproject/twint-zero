package main

import (
	"github.com/RSMCTECH/twint-zero/Core"
	"github.com/RSMCTECH/twint-zero/InputParser"
)

func main() {
	Arguments := InputParser.InputParser()
	Core.Main(&(Arguments.Query), &(Arguments.Instance), &(Arguments.Format))
}
