package interpreter

import (
	"bufio"
	"fmt"
	"github.com/paw1a/golox/internal/lexing"
	"github.com/paw1a/golox/internal/parsing"
	"io/ioutil"
	"os"
	"strings"
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

	sourceBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("can't read file %s: %v", filename, err)
		return
	}

	run(string(sourceBytes))

	if HasError {
		os.Exit(65)
	}
}

func runPrompt() {
	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("golox >>> ")

		line, err := in.ReadString('\n')
		if err != nil || strings.TrimSpace(line) == "exit" {
			break
		}

		run(line)
		HasError = false
	}
}

func run(source string) {
	lexer := lexing.NewLexer(source)
	lexer.ScanTokens()

	if len(lexer.Errors) == 0 {
		//for _, token := range lexer.Tokens {
		//	fmt.Printf("%s\n", token.String())
		//}
		parser := parsing.NewParser(lexer.Tokens)
		expr := parser.Parse()
		fmt.Printf("%s\n", expr.Print())
	} else {
		for _, err := range lexer.Errors {
			fmt.Printf("%s\n", err.Error())
		}
		HasError = true
	}
}

func PrintError(line int, message string) {
	fmt.Printf("line %d | error: %s", line, message)
	HasError = true
}
