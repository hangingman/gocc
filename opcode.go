package main

type Opcode int

const (
	MOVB Opcode = iota
	MOVW
	MOVL
	MOVQ
	ADDL
	SUBL
	SUBQ
	IMUL
	IDIV
	CLTD
	XORL
	PUSH
	POP
	CALL
	LEAVE
	RET
)

func mov(t CType) Opcode {
	switch t.Bytes() {
	case 1:
		return MOVB
	case 2:
		return MOVW
	case 4:
		return MOVL
	default:
		return MOVQ
	}
}

func (c Opcode) String() string {
	switch c {
	case MOVB:
		return "movb"
	case MOVW:
		return "movw"
	case MOVL:
		return "movl"
	case MOVQ:
		return "movq"
	case ADDL:
		return "addl"
	case SUBL:
		return "subl"
	case SUBQ:
		return "subq"
	case IMUL:
		return "imul"
	case IDIV:
		return "idiv"
	case CLTD:
		return "cltd"
	case XORL:
		return "xorl"
	case PUSH:
		return "push"
	case POP:
		return "pop "
	case CALL:
		return "call"
	case LEAVE:
		return "leave"
	case RET:
		return "ret"
	default:
		panic("undefined code")
	}
}
