package main

type Kind int

const (
	AST_VarDef Kind = iota
	AST_FuncDef
	AST_FuncArg
	AST_Ident
	// expr
	AST_BinaryExpr
	AST_CondExpr
	AST_UnaryExpr
	AST_AssignExpr
	AST_FuncCall
	AST_IntVal
	AST_CharVal
	// stmt
	AST_BlockStmt
	AST_ReturnStmt
	AST_ExprStmt
)

type CType int

const (
	C_int CType = iota
	C_void
	C_char
	C_float
	C_long
	C_short
	C_double
)

// func (t Type) Size() int {
// 	switch t {
// 	case C_int:
// 		return 4
// 	case C_Char:
// 		return 1
// 	default:
// 		panic("unimplemented type size")
// 	}
// }

func (t CType) String() string {
	switch t {
	case C_int:
		return "int"
	case C_void:
		return "void"
	case C_char:
		return "char"
	case C_float:
		return "float"
	case C_long:
		return "long"
	case C_short:
		return "short"
	case C_double:
		return "double"
	default:
		panic("undefined Type")
	}
}

type (
	Node interface {
		Kind() Kind
	}
)

type (
	Ident struct {
		Token *Token
	}

	VarDef struct {
		Type CType
		Name string
		Init *Expr
	}

	FuncDef struct {
		Type  CType
		Name  string
		Args  []FuncArg
		Block BlockStmt
	}

	FuncArg struct {
		Type CType
		Name *Token
	}
)

type (
	Expr interface {
		Node
		expr()
	}

	BinaryExpr struct {
		X  Expr
		Op *Token
		Y  Expr
	}

	CondExpr struct {
		Cond Expr
		L    Expr
		R    Expr
	}

	UnaryExpr struct {
		Op   *Token
		Expr Expr
	}

	AssignExpr struct {
		L  Expr
		Op *Token
		R  Expr
	}

	IntVal struct {
		Token *Token
	}

	CharVal struct {
		Token *Token
	}

	FuncCall struct {
		Ident Ident
		Args  []Expr
	}
)

type (
	Stmt interface {
		Node
		stmt()
	}

	BlockStmt struct {
		Nodes []Node
	}

	ReturnStmt struct {
		Expr Expr
	}

	ExprStmt struct {
		Expr Expr
	}
)

func (VarDef) Kind() Kind     { return AST_VarDef }
func (FuncDef) Kind() Kind    { return AST_FuncDef }
func (FuncArg) Kind() Kind    { return AST_FuncArg }
func (Ident) Kind() Kind      { return AST_Ident }
func (BinaryExpr) Kind() Kind { return AST_BinaryExpr }
func (CondExpr) Kind() Kind   { return AST_CondExpr }
func (UnaryExpr) Kind() Kind  { return AST_UnaryExpr }
func (AssignExpr) Kind() Kind { return AST_AssignExpr }
func (FuncCall) Kind() Kind   { return AST_FuncCall }
func (IntVal) Kind() Kind     { return AST_IntVal }
func (CharVal) Kind() Kind    { return AST_CharVal }
func (BlockStmt) Kind() Kind  { return AST_BlockStmt }
func (ReturnStmt) Kind() Kind { return AST_ReturnStmt }
func (ExprStmt) Kind() Kind   { return AST_ExprStmt }

func (Ident) expr()      {}
func (BinaryExpr) expr() {}
func (CondExpr) expr()   {}
func (UnaryExpr) expr()  {}
func (AssignExpr) expr() {}
func (FuncCall) expr()   {}
func (IntVal) expr()     {}
func (CharVal) expr()    {}

func (BlockStmt) stmt()  {}
func (ReturnStmt) stmt() {}
func (ExprStmt) stmt()   {}
