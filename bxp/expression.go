package bxp

import (
	"fmt"
	"strconv"
	"strings"
)

type Operation string

const (
	AND  Operation = "AND"
	NAND           = "NAND"
	OR             = "OR"
	NOR            = "NOR"
	XOR            = "XOR"
	XNOR           = "XNOR"
)

func FromString(str string) Operation {
	switch str {
	case "+", "|", "or", "OR":
		return OR
	case "!+", "!|", "nor", "NOR":
		return NOR
	case "*", "&", "and", "AND":
		return AND
	case "!*", "!&", "nand", "NAND":
		return NAND
	case "^", "xor", "XOR":
		return XOR
	case "!^", "xnor", "XNOR":
		return XNOR
	default:
		return Operation(str)
	}
}

var id = 0

type Expression struct {
	id        int
	first     func() bool
	eFirst    *Expression
	second    func() bool
	eSecond   *Expression
	Operation Operation
}

func newExpression(
	first func() bool,
	eFirst *Expression,
	second func() bool,
	eSecond *Expression,
	operation Operation,
) *Expression {
	e := Expression{
		id:        id,
		first:     first,
		eFirst:    eFirst,
		second:    second,
		eSecond:   eSecond,
		Operation: operation,
	}
	id++
	return &e
}

func (e *Expression) First() bool {
	if e.first != nil {
		return e.first()
	}
	if e.eFirst != nil {
		return e.eFirst.Compute()
	}
	return false
}

func (e *Expression) Second() bool {
	if e.second != nil {
		return e.second()
	}
	if e.eSecond != nil {
		return e.eSecond.Compute()
	}
	return false
}

func (e *Expression) ToString() string {
	return e.toString(0)
}

func (e *Expression) toString(level int) string {
	return "EXPR(" + strconv.Itoa(e.id) + ")\n" +
		strings.Repeat("\t", level+1) + "FIRST " + e.firstString(level) + "\n" +
		strings.Repeat("\t", level+1) + "OP " + string(e.Operation) + "\n" +
		strings.Repeat("\t", level+1) + "SECOND " + e.secondString(level)
}

func (e *Expression) hasFirst() bool {
	return e.first != nil || e.eFirst != nil
}

func (e *Expression) hasSecond() bool {
	return e.second != nil || e.eSecond != nil
}

func (e *Expression) hasOperation() bool {
	return len(e.Operation) > 0
}

func (e *Expression) firstString(level int) string {
	if e.eFirst != nil {
		return e.eFirst.toString(level + 1)
	}
	return fmt.Sprintf("%p", e.first)
}

func (e *Expression) secondString(level int) string {
	if e.eSecond != nil {
		return e.eSecond.toString(level + 1)
	}
	return fmt.Sprintf("%p", e.second)
}

func And(first func() bool, second func() bool) *Expression {
	return NewBuilder().First(first, nil).Second(second, nil).Operation(AND).Build()
}

func EAnd(first *Expression, second *Expression) *Expression {
	return NewBuilder().First(nil, first).Second(nil, second).Operation(AND).Build()
}

func Nand(first func() bool, second func() bool) *Expression {
	return NewBuilder().First(first, nil).Second(second, nil).Operation(NAND).Build()
}

func Or(first func() bool, second func() bool) *Expression {
	return NewBuilder().First(first, nil).Second(second, nil).Operation(OR).Build()

}

func Nor(first func() bool, second func() bool) *Expression {
	return NewBuilder().First(first, nil).Second(second, nil).Operation(NOR).Build()
}

func (e *Expression) Negate() {
	switch e.Operation {
	case AND:
		e.Operation = NAND
		return
	case NAND:
		e.Operation = AND
		return
	case OR:
		e.Operation = NOR
		return
	case NOR:
		e.Operation = OR
		return
	case XOR:
		e.Operation = XNOR
		return
	case XNOR:
		e.Operation = XOR
		return
	}
}

func (e *Expression) Compute() bool {
	//fmt.Printf("Computing exp( %d ) - %t %s %t\n", e.id, e.First(), e.Operation, e.Second())
	switch e.Operation {
	case AND:
		return e.First() && e.Second()
	case NAND:
		return !(e.First() && e.Second())
	case OR:
		return e.First() || e.Second()
	case NOR:
		return !(e.First() || e.Second())
	case XOR:
		return e.First() != e.Second()
	case XNOR:
		return e.First() == e.Second()
	default:
		fmt.Println("NOOP")
		return false
	}
}

type Builder struct {
	expression *Expression
}

func NewBuilder() *Builder {
	return &Builder{
		newExpression(nil, nil, nil, nil, ""),
	}
}

func (b *Builder) Full() bool {
	return b.expression.hasFirst() && b.expression.hasSecond() && b.expression.hasOperation()
}

func (b *Builder) First(first func() bool, eFirst *Expression) *Builder {
	//fmt.Printf("(%d) Setting first\n", b.expression.id)
	b.expression.first = first
	b.expression.eFirst = eFirst
	return b
}

func (b *Builder) Second(second func() bool, eSecond *Expression) *Builder {
	//fmt.Printf("(%d) Setting second\n", b.expression.id)
	b.expression.second = second
	b.expression.eSecond = eSecond
	return b
}

func (b *Builder) Operation(operation Operation) *Builder {
	//fmt.Printf("(%d) Setting op %s\n", b.expression.id, operation)
	b.expression.Operation = operation
	return b
}

func (b *Builder) Build() *Expression {
	return b.expression
}
