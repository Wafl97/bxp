package bxp

import (
	"bool_interp/util"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Row struct {
	data []bool
}

func NewRow(width int, index int) Row {
	rowBitStr := fmt.Sprintf("%0"+strconv.Itoa(width-1)+"b", index)

	data := make([]bool, width)
	rowBits := strings.Split(rowBitStr, "")
	for i := range rowBits {
		data[i] = rowBits[i] == "1"
	}

	return Row{
		data,
	}
}

func (r *Row) ToString() string {
	s := ""
	for i := range r.data {
		if r.data[i] {
			s += fmt.Sprintf(" %s |", util.Green("1"))
		} else {
			s += fmt.Sprintf(" %s |", util.Red("0"))
		}
	}
	return fmt.Sprintf("|%s \n", s)
}

type TruthTable struct {
	boolExp    string
	labels     []string
	table      []Row
	rowMap     map[string]bool
	expression *Expression
}

var ops map[string]bool = map[string]bool{
	"+": true, "|": true, // OR
	"*": true, "&": true, // AND
	"!": true, "'": true, // NOT
	"^": true, // xor
	"(": true, ")": true,
}

func NewTruthTableFromFile(fileName string) *TruthTable {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return NewTruthTable(string(bytes))
}

func NewTruthTable(boolExp string) *TruthTable {
	// get labels
	charSet := map[string]bool{}
	var labels []string

	chars := strings.Split(boolExp, "")
	for i := range chars {
		_, contained := ops[chars[i]]
		alreadyExist := charSet[chars[i]]
		if !contained && !alreadyExist {
			charSet[chars[i]] = true
			labels = append(labels, chars[i])
		}
	}
	labels = append(labels, "Q")

	rowMap := map[string]bool{}

	// make bool matrix
	// 2^labels
	table := make([]Row, int(math.Pow(2, float64(len(labels)-1))))
	for i := range table {
		table[i] = NewRow(len(labels), i)
	}

	tt := &TruthTable{
		boolExp: boolExp,
		labels:  labels,
		table:   table,
		rowMap:  rowMap,
	}

	tt.parseExpression()

	return tt
}

func (tt *TruthTable) getMapping(char string) bool {
	//fmt.Printf("Getting %s=%t ", char, tt.rowMap[char])
	return tt.rowMap[char]
}

func (tt *TruthTable) parseExpression() {
	chars := strings.Split(tt.boolExp, "")
	e, _ := tt._parseExpression(chars, 0)
	tt.expression = e
	//fmt.Println(e.ToString(0))
}

func (tt *TruthTable) _parseExpression(chars []string, index int) (*Expression, int) {
	eb := NewBuilder()

	e, ef, i := tt.handleInput(chars, index)
	index = i // skip to operation
	eb.First(e, ef)

	op, i := tt.handleOperation(chars, index)
	index = i
	eb.Operation(op)

	e, ef, i = tt.handleInput(chars, index)
	index = i // skip to next
	eb.Second(e, ef)

	return eb.Build(), index
}

func (tt *TruthTable) handleInput(chars []string, index int) (func() bool, *Expression, int) {
	char := chars[index]
	//fmt.Printf("Handeling char[%d] = %s\n", index, char)
	if char == "(" {
		//fmt.Println("EXPR")
		e, i := tt._parseExpression(chars, index+1)
		return nil, e, i + 1
	}
	if char == "!" {
		if chars[index+1] == "(" {
			//fmt.Println("! EXPR")
			e, i := tt._parseExpression(chars, index+2)
			e.Negate()
			return nil, e, i + 1
		}
		//fmt.Println("!", chars[index+1])
		return func() bool { return !tt.getMapping(chars[index+1]) }, nil, index + 2
	}
	//fmt.Println(chars[index])
	return func() bool { return tt.getMapping(chars[index]) }, nil, index + 1
}

func (tt *TruthTable) handleOperation(chars []string, index int) (Operation, int) {
	if chars[index] == "!" {
		return FromString(chars[index] + chars[index+1]), index + 2
	}
	return FromString(chars[index]), index + 1
}

func (tt *TruthTable) Compute() {
	for i := range tt.table {
		tt.copyRowIntoMap(i)
		tt.table[i].data[len(tt.table[i].data)-1] = tt.expression.Compute()
	}
}

func (tt *TruthTable) copyRowIntoMap(rowIndex int) {
	for i := range tt.labels {
		tt.rowMap[tt.labels[i]] = tt.table[rowIndex].data[i]
	}
}

func (tt *TruthTable) GetTrueRowIndices() []int {
	var indices []int
	for i := range tt.table {
		if tt.table[i].data[len(tt.table[i].data)-1] {
			indices = append(indices, i)
		}
	}
	return indices
}

func (tt *TruthTable) GetRow(index int) Row {
	return tt.table[index]
}

func (tt *TruthTable) GetExpression() *Expression {
	return tt.expression
}

func (tt *TruthTable) ToString() string {
	rs := ""
	for i := range tt.table {
		rs += fmt.Sprintf("| %3d %s", i, tt.table[i].ToString())
	}

	return fmt.Sprintf("%s\n| row | %s |\n%s\n%s",
		util.Blue(fmt.Sprintf("Expression: %s", tt.boolExp)),
		strings.Join(tt.labels, " | "),
		fmt.Sprintf("|=====|%s", strings.Repeat("===|", len(tt.labels))),
		rs,
	)
}
