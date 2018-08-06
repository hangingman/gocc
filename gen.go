package main

import "fmt"

type Map map[string]int

type Gen struct {
	s   string
	pos int
	m   Map
}

func NewGen() *Gen {
	return &Gen{s: "", pos: 0, m: Map{}}
}

type Code int

const (
	ADDL Code = iota
	SUBL
	MOVL
	IMULL
	IDIVL
	CLTD
	PUSHL
	POPL
	LEAVE
	RET
)

type Reg int

const (
	EAX Reg = iota
	EBX
	ECX
	EDX

	EBP
	ESP
)

type Operand interface {
	Str() string
}

func (gen *Gen) add(n string, p int) {
	gen.m[n] = p
}

func (gen *Gen) lookup(n string) (int, bool) {
	for k, v := range gen.m {
		if k == n {
			return v, true
		}
	}
	return 0, false
}

func (r Reg) Str() string    { return r.String() }
func (i IntVal) Str() string { return "$" + string(i.Token.Str) }

func (gen *Gen) emit(c Code, ops ...Operand) {
	gen.s += "\t" + c.String()
	for i, v := range ops {
		if i != 0 {
			gen.s += ","
		}
		gen.s += "\t" + v.Str()
	}
	gen.s += "\n"
}

func (gen *Gen) global(n string) {
	gen.s += ".global " + n + "\n"
}

func (gen *Gen) prologue() {
	gen.emit(PUSHL, EBP)
	gen.emit(MOVL, ESP, EBP)
}

func (gen *Gen) epilogue() {
	gen.emit(LEAVE)
	gen.emit(RET)
}

func (gen *Gen) emitFuncDef(n string) {
	gen.s += n + ":\n"
}

func (gen *Gen) generate(n Node) {
	switch v := n.(type) {
	case FuncDef:
		if v.Name == "main" {
			gen.global("_main")
			gen.emitFuncDef("_main")
		} else {
			gen.emitFuncDef(v.Name)
		}
		gen.prologue()

		for _, node := range v.Block.Nodes {
			gen.generate(node)
		}

		gen.epilogue()

	case Expr:
		gen.expr(v)
	case VarDef:
		gen.varDef(v)
	default:
		panic("unimplemented")
	}
}

func (gen *Gen) expr(e Expr) {
	switch v := e.(type) {
	case BinaryExpr:
		gen.binary(v)
	case Ident:
		if pos, ok := gen.lookup(v.Token.String()); ok {
			gen.s += fmt.Sprintf("\tmovl\t%d(%%ebp), %%eax\n", -pos)
		} else {
			panic("ident is not defined")
		}
	case IntVal:
		gen.emit(MOVL, v, EAX)
	}
}

func (gen *Gen) binary(e BinaryExpr) {
	gen.expr(e.X)
	gen.emit(PUSHL, EAX)

	gen.expr(e.Y)
	gen.emit(MOVL, EAX, EBX)

	gen.emit(POPL, EAX)

	var c Code
	if e.Op.Kind == ADD {
		c = ADDL
	} else if e.Op.Kind == SUB {
		c = SUBL
	} else if e.Op.Kind == MUL {
		c = IMULL
	} else if e.Op.Kind == DIV || e.Op.Kind == REM {
		c = IDIVL
	} else {
		panic("unimplemented")
	}

	if c == IDIVL {
		gen.emit(CLTD)
		gen.emit(IDIVL, EBX)
		if e.Op.Kind == REM {
			gen.emit(MOVL, EDX, EAX)
		}
	} else {
		gen.emit(c, EBX, EAX)
	}
}

func (gen *Gen) varDef(n VarDef) {
	if n.Init != nil {
		gen.expr(*n.Init)
	}
	gen.pos += n.Type.Size()
	gen.add(n.Name, gen.pos)
	gen.s += fmt.Sprintf("\tsubl\t$%d, %%esp\n", n.Type.Size())
	gen.s += fmt.Sprintf("\tmovl\t%%eax, %d(%%ebp)\n", -gen.pos)
}

func (c Code) String() string {
	switch c {
	case ADDL:
		return "addl"
	case SUBL:
		return "subl"
	case MOVL:
		return "movl"
	case IMULL:
		return "imull"
	case IDIVL:
		return "idivl"
	case CLTD:
		return "cltd"
	case PUSHL:
		return "pushl"
	case POPL:
		return "popl"
	case LEAVE:
		return "leave"
	case RET:
		return "ret"
	default:
		panic("undefined code")
	}
}

func (r Reg) String() string {
	switch r {
	case EAX:
		return "%eax"
	case EBX:
		return "%ebx"
	case ECX:
		return "%ecx"
	case EDX:
		return "%edx"
	case EBP:
		return "%ebp"
	case ESP:
		return "%esp"
	default:
		panic("undefined reg")
	}
}
