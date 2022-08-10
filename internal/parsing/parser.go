package parsing

import (
	"bytes"
	"fmt"
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
	"strconv"
	"strings"
)

type Parser struct {
	tokens  []lexing.Token
	current int

	Errors []error
	lines  []string

	isLoopScope bool
	isFuncScope bool
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt

	for !p.isEof() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() ast.Stmt {
	defer p.parseRecoverFunc()

	if p.match(lexing.Var) {
		p.advance()
		return p.varDeclaration()
	}

	if p.match(lexing.Fun) {
		p.advance()
		return p.funDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() ast.Stmt {
	varName := p.requireToken(lexing.Identifier, "variable name expected")

	var initializer ast.Expr
	if p.peek().TokenType == lexing.Equal {
		p.advance()
		initializer = p.expression()
	}

	p.requireToken(lexing.Semicolon, "';' expected")
	return ast.VarDeclarationStmt{
		Name:        varName,
		Initializer: initializer,
	}
}

func (p *Parser) funDeclaration() ast.Stmt {
	funcName := p.requireToken(lexing.Identifier, "function name expected")
	p.requireToken(lexing.LeftParen, "function declaration expect '('")

	parameters := make([]lexing.Token, 0)
	if !p.match(lexing.RightParen) {
		parameters = append(parameters, p.requireToken(lexing.Identifier,
			"function declaration expect identifier as param name"))
		for p.match(lexing.Comma) {
			p.advance()
			if len(parameters) >= 255 {
				p.parseError(p.peek(), "declared more than 255 parameters")
				break
			}
			parameters = append(parameters, p.requireToken(lexing.Identifier,
				"function declaration expect identifier as param name"))
		}
	}
	p.requireToken(lexing.RightParen, "function declaration expect ')'")

	p.requireToken(lexing.LeftBrace, "expect '{' before function body")

	innerFunc := p.isFuncScope
	p.isFuncScope = true
	defer func() {
		if !innerFunc {
			p.isFuncScope = false
		}
	}()

	statement := p.blockStatement()

	return ast.FunDeclarationStmt{
		Name:      funcName,
		Params:    parameters,
		Statement: statement.(ast.BlockStmt),
	}
}

func (p *Parser) statement() ast.Stmt {
	switch {
	case p.match(lexing.LeftBrace):
		p.advance()
		return p.blockStatement()
	case p.match(lexing.If):
		p.advance()
		return p.ifStatement()
	case p.match(lexing.While):
		p.advance()
		return p.whileStatement()
	case p.match(lexing.For):
		p.advance()
		return p.forStatement()
	case p.match(lexing.Break):
		if p.isLoopScope {
			p.advance()
			return p.breakStatement()
		} else {
			p.parseError(p.peek(), "break statement not within loop")
		}
	case p.match(lexing.Continue):
		if p.isLoopScope {
			p.advance()
			return p.continueStatement()
		} else {
			p.parseError(p.peek(), "continue statement not within loop")
		}
	case p.match(lexing.Return):
		if p.isFuncScope {
			return p.returnStatement(p.advance())
		} else {
			p.parseError(p.peek(), "return statement not within func body")
		}
	}

	return p.expressionStatement()
}

func (p *Parser) returnStatement(returnToken lexing.Token) ast.Stmt {
	var returnExpr ast.Expr
	if !p.match(lexing.Semicolon) {
		returnExpr = p.expression()
	}
	p.requireToken(lexing.Semicolon, "expect ';' after return keyword")

	return ast.ReturnStmt{
		ReturnToken: returnToken,
		Expr:        returnExpr,
	}
}

func (p *Parser) breakStatement() ast.Stmt {
	p.requireToken(lexing.Semicolon, "expect ';' after break statement")
	return ast.BreakStmt{}
}

func (p *Parser) continueStatement() ast.Stmt {
	p.requireToken(lexing.Semicolon, "expect ';' after continue statement")
	return ast.ContinueStmt{}
}

func (p *Parser) forStatement() ast.Stmt {
	p.requireToken(lexing.LeftParen, "for statement expect '('")

	var initializerStmt ast.Stmt
	switch {
	case p.match(lexing.Var):
		p.advance()
		initializerStmt = p.varDeclaration()
	case p.match(lexing.Semicolon):
		p.advance()
	default:
		initializerStmt = p.expressionStatement()
	}

	var conditionExpr ast.Expr
	if !p.match(lexing.Semicolon) {
		conditionExpr = p.expression()
	} else {
		conditionExpr = ast.LiteralExpr{LiteralValue: true}
	}
	p.requireToken(lexing.Semicolon,
		"for statement expect ';' between condition and increment expressions")

	var incrementExpr ast.Expr
	if !p.match(lexing.RightParen) {
		incrementExpr = p.expression()
	}
	p.requireToken(lexing.RightParen, "for statement expect ')'")

	innerLoop := p.isLoopScope
	p.isLoopScope = true
	defer func() {
		if !innerLoop {
			p.isLoopScope = false
		}
	}()

	statement := p.statement()

	return ast.ForStmt{
		InitializerStmt: initializerStmt,
		ConditionExpr:   conditionExpr,
		IncrementExpr:   incrementExpr,
		Statement:       statement,
	}
}

func (p *Parser) whileStatement() ast.Stmt {
	p.requireToken(lexing.LeftParen, "while statement expect '(' before condition")
	conditionExpr := p.expression()
	p.requireToken(lexing.RightParen, "while statement expect ')' after condition")

	innerLoop := p.isLoopScope
	p.isLoopScope = true
	defer func() {
		if !innerLoop {
			p.isLoopScope = false
		}
	}()

	statement := p.statement()

	return ast.ForStmt{
		ConditionExpr: conditionExpr,
		Statement:     statement,
	}
}

func (p *Parser) ifStatement() ast.Stmt {
	p.requireToken(lexing.LeftParen, "if statement expect '(' before condition")
	conditionExpr := p.expression()
	p.requireToken(lexing.RightParen, "if statement expect ')' after condition")
	ifStatement := p.statement()

	var elseStatement ast.Stmt
	if p.match(lexing.Else) {
		p.advance()
		elseStatement = p.statement()
	}

	return ast.IfStmt{
		ConditionExpr: conditionExpr,
		IfStatement:   ifStatement,
		ElseStatement: elseStatement,
	}
}

func (p *Parser) blockStatement() ast.Stmt {
	var stmts []ast.Stmt

	for !p.match(lexing.RightBrace) && !p.isEof() {
		stmts = append(stmts, p.declaration())
	}

	p.requireToken(lexing.RightBrace, "'}' end of block expected")

	return ast.BlockStmt{Stmts: stmts}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.requireToken(lexing.Semicolon, "';' expected")
	return ast.ExpressionStmt{Expr: expr}
}

func (p *Parser) expression() ast.Expr {
	return p.comma()
}

func (p *Parser) lambda() ast.Expr {
	p.requireToken(lexing.LeftParen, "lambda declaration expect '('")

	parameters := make([]lexing.Token, 0)
	if !p.match(lexing.RightParen) {
		parameters = append(parameters, p.requireToken(lexing.Identifier,
			"lambda declaration expect identifier as param name"))
		for p.match(lexing.Comma) {
			p.advance()
			if len(parameters) >= 255 {
				p.parseError(p.peek(), "declared more than 255 parameters")
				break
			}
			parameters = append(parameters, p.requireToken(lexing.Identifier,
				"lambda declaration expect identifier as param name"))
		}
	}
	p.requireToken(lexing.RightParen, "lambda declaration expect ')'")

	p.requireToken(lexing.LeftBrace, "expect '{' before lambda body")

	innerFunc := p.isFuncScope
	p.isFuncScope = true
	defer func() {
		if !innerFunc {
			p.isFuncScope = false
		}
	}()

	statement := p.blockStatement()

	return ast.LambdaExpr{
		Params:    parameters,
		Statement: statement.(ast.BlockStmt),
	}
}

func (p *Parser) comma() ast.Expr {
	var expr ast.Expr

	expr = p.assignment()
	for p.match(lexing.Comma) {
		operator := p.advance()
		rightExpr := p.assignment()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) assignment() ast.Expr {
	if p.match(lexing.Fun) {
		p.advance()
		return p.lambda()
	}

	expr := p.logicalOr()

	if p.match(lexing.Equal) {
		equalToken := p.advance()
		value := p.assignment()

		switch expr.(type) {
		case ast.VariableExpr:
			return ast.AssignExpr{
				Variable:    expr.(ast.VariableExpr),
				Initializer: value,
			}
		case ast.IndexExpr:
			return ast.AssignExpr{
				Variable:    expr.(ast.IndexExpr),
				Initializer: value,
			}
		}

		p.parseError(equalToken, "invalid assignment target")
	}

	if p.match(lexing.Question) {
		p.advance()
		trueValue := p.logicalOr()
		p.requireToken(lexing.Colon, "ternary operator expect ':'")
		falseValue := p.logicalOr()
		return ast.TernaryExpr{
			Condition: expr,
			TrueExpr:  trueValue,
			FalseExpr: falseValue,
		}
	}

	return expr
}

func (p *Parser) logicalOr() ast.Expr {
	var expr ast.Expr

	expr = p.logicalAnd()
	for p.match(lexing.Or) {
		operator := p.advance()
		rightExpr := p.logicalAnd()
		expr = ast.LogicalExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) logicalAnd() ast.Expr {
	var expr ast.Expr

	expr = p.equality()
	for p.match(lexing.And) {
		operator := p.advance()
		rightExpr := p.equality()
		expr = ast.LogicalExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) equality() ast.Expr {
	var expr ast.Expr

	expr = p.comparison()
	for p.match(lexing.BangEqual, lexing.EqualEqual) {
		operator := p.advance()
		rightExpr := p.comparison()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	var expr ast.Expr

	expr = p.term()
	for p.match(lexing.Less, lexing.LessEqual, lexing.Greater, lexing.GreaterEqual) {
		operator := p.advance()
		rightExpr := p.term()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	var expr ast.Expr

	expr = p.factor()
	for p.match(lexing.Minus, lexing.Plus) {
		operator := p.advance()
		rightExpr := p.factor()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	var expr ast.Expr

	expr = p.unary()
	for p.match(lexing.Star, lexing.Slash) {
		operator := p.advance()
		rightExpr := p.unary()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	for p.match(lexing.Bang, lexing.Minus) {
		operator := p.advance()
		rightExpr := p.unary()
		return ast.UnaryExpr{
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return p.call()
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

	if p.match(lexing.LeftParen) {
		for p.match(lexing.LeftParen) {
			p.advance()
			expr = p.callArguments(expr)
		}
		return expr
	}

	if p.match(lexing.LeftBracket) {
		for p.match(lexing.LeftBracket) {
			p.advance()
			expr = p.arrayIndex(expr)
		}
		return expr
	}

	return expr
}

func (p *Parser) callArguments(callee ast.Expr) ast.Expr {
	arguments := make([]ast.Expr, 0)

	if !p.match(lexing.RightParen) {
		arguments = append(arguments, p.assignment())
		for p.match(lexing.Comma) {
			p.advance()
			if len(arguments) >= 255 {
				p.parseError(p.peek(), "passed more than 255 arguments")
				break
			}
			arguments = append(arguments, p.assignment())
		}
	}

	paren := p.requireToken(lexing.RightParen, "function call expect ')'")
	return ast.CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (p *Parser) arrayIndex(array ast.Expr) ast.Expr {
	indexExpr := p.expression()
	bracket := p.requireToken(lexing.RightBracket, "array index expression expect ']'")
	return ast.IndexExpr{
		Array:     array,
		Bracket:   bracket,
		IndexExpr: indexExpr,
	}
}

func (p *Parser) primary() ast.Expr {
	switch {
	case p.match(lexing.False):
		p.advance()
		return ast.LiteralExpr{LiteralValue: false}
	case p.match(lexing.True):
		p.advance()
		return ast.LiteralExpr{LiteralValue: true}
	case p.match(lexing.Nil):
		p.advance()
		return ast.LiteralExpr{LiteralValue: nil}
	case p.match(lexing.Number, lexing.String):
		astNode := ast.LiteralExpr{LiteralValue: p.peek().Literal}
		p.advance()
		return astNode
	case p.match(lexing.LeftParen):
		p.advance()
		expr := p.expression()
		p.requireToken(lexing.RightParen, "expect ')' token after expression")
		return ast.GroupingExpr{Expr: expr}
	case p.match(lexing.LeftBracket):
		p.advance()
		return p.arrayElements()
	case p.match(lexing.Identifier):
		return ast.VariableExpr{Name: p.advance()}
	}

	p.parseError(p.peek(), "expect expression")
	return nil
}

func (p *Parser) arrayElements() ast.Expr {
	elements := make([]ast.Expr, 0)

	if !p.match(lexing.RightBracket) {
		elements = append(elements, p.assignment())
		for p.match(lexing.Comma) {
			p.advance()
			elements = append(elements, p.assignment())
		}
	}

	p.requireToken(lexing.RightBracket, "array initializer expect ']'")
	return ast.ArrayExpr{
		Elements: elements,
	}
}

func (p *Parser) requireToken(tokenType lexing.TokenType, message string) lexing.Token {
	if p.match(tokenType) {
		return p.advance()
	} else {
		p.parseError(p.peek(), message)
		return lexing.Token{}
	}
}

func (p *Parser) parseError(token lexing.Token, message string) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("[ %d:%d ]: error: %s\n",
		token.Line, token.Position, message))

	lineStr := strconv.Itoa(token.Line)
	buffer.WriteString(fmt.Sprintf("      %d |         %s\n", token.Line, p.lines[token.Line-1]))
	buffer.WriteString(fmt.Sprintf("      "))
	buffer.WriteString(strings.Repeat(" ", len(lineStr)))
	buffer.WriteString(" |         ")
	buffer.WriteString(fmt.Sprintf("%s^", strings.Repeat(" ", token.Position)))

	if len(token.Lexeme) > 0 {
		buffer.WriteString(fmt.Sprintf("%s\n", strings.Repeat("~", len(token.Lexeme)-1)))
	}

	panic(buffer.String())
}

func (p *Parser) parseRecoverFunc() {
	if err := recover(); err != nil {
		p.Errors = append(p.Errors, fmt.Errorf("%v", err))
		p.synchronize()
	}
}

func (p *Parser) synchronize() {
	for !p.isEof() {
		if p.match(lexing.Semicolon, lexing.Class, lexing.Fun,
			lexing.For, lexing.If, lexing.While,
			lexing.Return, lexing.Var) {
			p.advance()
			return
		}
		p.advance()
	}
}

func NewParser(tokens []lexing.Token, lines []string) *Parser {
	return &Parser{
		tokens: tokens,
		lines:  lines,
	}
}
