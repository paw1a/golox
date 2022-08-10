package interpreter

import (
	"bufio"
	"fmt"
	"github.com/paw1a/golox/internal/lexing"
	"github.com/paw1a/golox/internal/parsing"
	"github.com/paw1a/golox/internal/resolving"
	"github.com/paw1a/golox/internal/runtime"
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
}

func runPrompt() {
	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("golox >>> ")

		line, err := in.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil || line == "exit" {
			break
		}

		if !strings.HasSuffix(line, ";") && !strings.HasSuffix(line, "}") {
			line += ";"
		}

		run(line)
	}
}

func run(source string) {
	lexer := lexing.NewLexer(source)
	lexer.ScanTokens()

	if len(lexer.Errors) != 0 {
		for _, err := range lexer.Errors {
			fmt.Printf("%s\n", err.Error())
		}
		HasError = true
		return
	}

	//for _, token := range lexer.Tokens {
	//	fmt.Printf("%d %d\n", token.Line, token.Position)
	//}

	parser := parsing.NewParser(lexer.Tokens, lexer.Lines)
	statements := parser.Parse()

	if len(parser.Errors) != 0 {
		for _, err := range parser.Errors {
			fmt.Printf("%s\n", err.Error())
		}
		HasError = true
		return
	}

	inter := runtime.NewInterpreter()
	resolver := resolving.NewResolver(inter)

	defer errorRecovery()
	for _, stmt := range statements {
		resolver.ResolveStmt(stmt)
	}

	defer errorRecovery()
	for _, stmt := range statements {
		inter.Execute(stmt)
	}
}

func errorRecovery() {
	if err := recover(); err != nil {
		fmt.Printf("%v\n", err)
	}
}
