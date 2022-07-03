package interpreter

import (
	"fmt"
	"io/ioutil"
	"os"
)

var HasError = false

func Run() {
	if len(os.Args) > 2 {
		fmt.Printf("usage: golox [source code filename]")
		return
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("invalid source code filename: %s", filename)
		return
	}

	textBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("can't read file %s: %v", filename, err)
		return
	}

	run(string(textBytes))

	if HasError {
		os.Exit(65)
	}
}

func runPrompt() {
	for {
		fmt.Printf("golox >>> ")

		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			break
		}

		run(line)
		HasError = false
	}
}

func run(sourceCode string) {

}
