package main

import (
	"fmt"
)

type Parser struct {
	lexer *Lexer
	token *Token
	stack *Stack
}

func NewParser(source []byte) *Parser {
	p := &Parser{lexer: NewLexer(source), token: NewToken(), stack: &Stack{}}
	p.next()
	return p
}

func (p *Parser) match(t TokenKind) bool {
	return p.token.Kind == t
}

func (p *Parser) matchs(ts []TokenKind) bool {
	for _, v := range ts {
		if p.match(v) {
			return true
		}
	}
	return false
}

// push current token and position to stack
func (p *Parser) push() {
	pos := p.lexer.Pos()
	p.stack.push(*p.token, pos)
}

// pop last pushed token and position
func (p *Parser) pop() {
	t, pos := p.stack.pop()
	p.token = t
	p.lexer.Reset(pos)
}

// just remove last pushed token and position from stack
func (p *Parser) remove() {
	p.stack.pop()
}

var (
	typeKeys = []TokenKind{
		INT,
		VOID,
		CHAR,
		FLOAT,
		LONG,
		SHORT,
		DOUBLE,
		STRUCT,
		UNION,
		SIGNED,
		UNSIGNED,
		STATIC,
		AUTO,
		EXTERN,
		REGISTER,
		CONST,
		VOLATILE,
	}

	unaryOps = []TokenKind{
		AND,
		MUL,
		ADD,
		SUB,
		TILDE,
		NOT,
	}

	assignOps = []TokenKind{
		ASSIGN,
		MUL_ASSIGN,
		DIV_ASSIGN,
		REM_ASSIGN,
		ADD_ASSIGN,
		SUB_ASSIGN,
		LEFT_ASSIGN,
		RIGHT_ASSIGN,
		AND_ASSIGN,
		OR_ASSIGN,
		XOR_ASSIGN,
	}

	storageSpecifiers = []TokenKind{
		AUTO,
		REGISTER,
		STATIC,
		EXTERN,
		TYPEDEF,
	}

	typeSpecifiers = []TokenKind{
		VOID,
		CHAR,
		SHORT,
		INT,
		SHORT,
		LONG,
		FLOAT,
		DOUBLE,
		SIGNED,
		UNSIGNED,
	}
)

func (p *Parser) assert(t TokenKind) {
	if !p.match(t) {
		str := fmt.Sprintf("expected token is '"+t.String()+"', but got '"+p.token.String()+"' at line %d column %d", p.lexer.Pos().Line, p.lexer.Pos().Column)
		panic(str)
	}
}

func (p *Parser) next() {
	p.token = p.lexer.Next()
}

func (p *Parser) isEnd() bool {
	return p.match(EOF)
}

func (p *Parser) parse() Node {
	if p.isType() {
		return p.readVarDef()
	} else {
		return p.expr()
	}
}

func (p *Parser) expr() Expr {
	return p.assignExpr()
}

func (p *Parser) assignExpr() Expr {
	p.push()

	var hasAssign bool
	for !p.match(SEMICOLON) {
		if p.match(ASSIGN) {
			hasAssign = true
			break
		}
		p.next()

		if p.match(EOF) {
			break
		}
	}

	p.pop()

	if hasAssign {
		L := p.unaryExpr()
		if !p.isAssignOp() {
			panic("expected assign operator")
		}
		op := p.token
		p.next()
		R := p.assignExpr()
		p.assert(SEMICOLON)
		p.next()

		n := AssignExpr{L: L, Op: op, R: R}

		// look up ident name was declared before
		// ident, _ := L.(Ident)
		// name := string(ident.Str)
		// obj := p.lookup(name)
		// if obj == nil {
		// 	panic(name + " is not declared")
		// }
		// obj.IsInit = true

		return n
	} else {
		return p.conditionalExpr()
	}
}

func (p *Parser) isAssignOp() bool {
	return p.matchs(assignOps)
}

func (p *Parser) conditionalExpr() Expr {
	e := p.logOrExpr()
	if p.match(QUE) {
		p.next()
		L := p.expr()
		p.assert(COLON)
		p.next()
		n := CondExpr{Cond: e, L: L, R: p.conditionalExpr()}
		return n
	}
	return e
}

func (p *Parser) logOrExpr() Expr {
	e := p.logAndExpr()
	return p.logOrExpr2(e)
}

func (p *Parser) logOrExpr2(e Expr) Expr {
	if p.match(LOR) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.logAndExpr()}
		return p.logOrExpr2(n)
	}
	return e
}

func (p *Parser) logAndExpr() Expr {
	e := p.incOrExpr()
	return p.logAndExpr2(e)
}

func (p *Parser) logAndExpr2(e Expr) Expr {
	if p.match(LAND) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.incOrExpr()}
		return p.logAndExpr2(n)
	}
	return e
}

func (p *Parser) incOrExpr() Expr {
	e := p.excOrExpr()
	return p.incOrExpr2(e)
}

func (p *Parser) incOrExpr2(e Expr) Expr {
	if p.match(OR) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.excOrExpr()}
		return p.incOrExpr2(n)
	}
	return e
}

func (p *Parser) excOrExpr() Expr {
	e := p.andExpr()
	return p.excOrExpr2(e)
}

func (p *Parser) excOrExpr2(e Expr) Expr {
	if p.match(XOR) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.andExpr()}
		return p.excOrExpr2(n)
	}
	return e
}

func (p *Parser) andExpr() Expr {
	e := p.eqExpr()
	return p.andExpr2(e)
}

func (p *Parser) andExpr2(e Expr) Expr {
	if p.match(AND) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.eqExpr()}
		return p.andExpr2(n)
	}
	return e
}

func (p *Parser) eqExpr() Expr {
	e := p.relExpr()
	return p.eqExpr2(e)
}

func (p *Parser) eqExpr2(e Expr) Expr {
	if p.match(EQ) || p.match(NE) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.relExpr()}
		return p.eqExpr2(n)
	}
	return e
}

func (p *Parser) relExpr() Expr {
	e := p.shiftExpr()
	return p.relExpr2(e)
}

func (p *Parser) relExpr2(e Expr) Expr {
	if p.match(LT) || p.match(GT) || p.match(LE) || p.match(GE) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.shiftExpr()}
		return p.relExpr2(n)
	}
	return e
}

func (p *Parser) shiftExpr() Expr {
	e := p.additiveExpr()
	return p.shiftExpr2(e)
}

func (p *Parser) shiftExpr2(e Expr) Expr {
	if p.match(LSHIFT) || p.match(RSHIFT) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.additiveExpr()}
		return p.shiftExpr2(n)
	}
	return e
}

func (p *Parser) additiveExpr() Expr {
	e := p.multiExpr()
	return p.additiveExpr2(e)
}

func (p *Parser) additiveExpr2(e Expr) Expr {
	if p.match(ADD) || p.match(SUB) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.multiExpr()}
		return p.additiveExpr2(n)
	}
	return e
}

func (p *Parser) multiExpr() Expr {
	e := p.castExpr()
	return p.multiExpr2(e)
}

func (p *Parser) multiExpr2(e Expr) Expr {
	if p.match(MUL) || p.match(DIV) || p.match(REM) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.castExpr()}
		return p.multiExpr2(n)
	}
	return e
}

func (p *Parser) castExpr() Expr {
	return p.unaryExpr()
}

func (p *Parser) unaryExpr() Expr {
	if p.match(INC) || p.match(DEC) {
		panic("unimplemented unaryExpr")
	} else if p.isUnaryOp() {
		op := p.token
		p.next()

		return UnaryExpr{Op: op, E: p.castExpr()}
	} else {
		return p.postfixExpr()
	}
}

func (p *Parser) isUnaryOp() bool {
	return p.matchs(unaryOps)
}

func (p *Parser) postfixExpr() Expr {
	n := p.primaryExpr()
	return p.postfixExpr2(n)
}

func (p *Parser) postfixExpr2(e Expr) Expr {
	if p.match(INC) {
		panic("postfix increment")
	} else if p.match(DEC) {
		panic("postfix decrement")
	} else if p.match(LPAREN) {
		switch e.(type) {
		// case Ident:
		// 	return p.readFuncCall(e)
		default:
			panic("unimplemented postfixExpr2")
		}
	} else if p.match(PERIOD) {
		panic("postfix .")
	} else if p.match(ARROW) {
		panic("postfix ->")
	} else {
		return e
	}
}

func (p *Parser) primaryExpr() Expr {
	switch {
	case p.match(IDENT):
		n := Ident{Token: p.token}
		p.next()
		return n
	case p.match(INT_CONST):
		n := IntVal{Token: p.token}
		p.next()
		return n
	case p.match(LPAREN):
		p.next()
		e := p.expr()
		p.assert(RPAREN)
		p.next()
		return e
	default:
		panic("primaryExpr")
	}
}

func (p *Parser) isType() bool {
	return p.matchs([]TokenKind{INT, CHAR, VOID, FLOAT, LONG, SHORT, DOUBLE})
}

func (p *Parser) readType() Type {
	var t Type
	switch p.token.Kind {
	case INT:
		t = Int_t
	case CHAR:
		t = Char_t
	case VOID:
		t = Void_t
	case FLOAT:
		t = Float_t
	case LONG:
		t = Long_t
	case SHORT:
		t = Short_t
	case DOUBLE:
		t = Double_t
	default:
		panic("readType")
	}
	p.next()
	return t
}

func (p *Parser) readVarDef() Node {
	t := p.readType()

	p.assert(IDENT)
	name := string(p.token.Str)
	p.next()

	var n Node
	if p.match(LBRACK) {
		panic("read subscript")
		// var s []Expr
		// for p.match(LBRACK) {
		// 	s = append(s, p.readSubscript())
		// }
		// arr := &ast.ArrayDef{Type: t, Name: name, Subscript: s}
		//
		// if p.match(ASSIGN) {
		// 	p.next()
		// 	arr.Init = p.readArrayInit()
		// 	// obj.IsInit = true
		// }
		// n = arr
	} else {
		v := VarDef{Type: t, Name: name}

		if p.match(ASSIGN) {
			p.next()
			e := p.assignExpr()
			v.Init = &e
		}
		n = v
	}

	p.assert(SEMICOLON)
	p.next()

	return n
}

func (p *Parser) readFuncDef() FuncDef {
	t := p.readType()

	p.assert(IDENT)
	name := string(p.token.Str)
	p.next()

	p.assert(LPAREN)
	p.next()

	args := p.readFuncArgs()

	p.assert(RPAREN)
	p.next()

	block := p.blockStmt()

	return FuncDef{Type: t, Name: name, Args: args, Block: block}
}

func (p *Parser) readFuncArgs() []FuncArg {
	var res []FuncArg
	for {
		res = append(res, p.readFuncArg())
		if !p.match(COMMA) {
			break
		}
		p.next()
	}
	return res
}

func (p *Parser) readFuncArg() FuncArg {
	var n FuncArg
	n.Type = p.readType()

	p.assert(IDENT)
	n.Name = p.token
	p.next()

	return n
}

func (p *Parser) blockStmt() BlockStmt {
	p.assert(LBRACE)
	p.next()
	n := BlockStmt{}

	for !p.match(RBRACE) {
		if p.isType() {
			d := p.readVarDef()
			n.Nodes = append(n.Nodes, d)
		} else {
			// - TODO:
			panic("stmt is not implemented")
			// stmt := p.stmt()
			// n.Nodes = append(n.Nodes, stmt)
		}
	}
	p.next()

	return n
}
