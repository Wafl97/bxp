package main

import (
	"bool_interp/bxp"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// TODO: cleanup this mess
func main() {
	filename := flag.String("f", "", "the file to read from")
	boolExpression := flag.String("b", "", "some boolean expression")
	showTable := flag.Bool("p", false, "print the truth table")
	showTrueRows := flag.Bool("q", false, "get indices all true rows")
	showModel := flag.Bool("m", false, "show the expression model")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	var tt *bxp.TruthTable

	if len(*filename) > 0 {
		tt = bxp.NewTruthTableFromFile(*filename)
	} else if len(*boolExpression) > 0 {
		tt = bxp.NewTruthTable(*boolExpression)
	} else {
		*showTable = true
		fmt.Println("Enter empty line to exit")
		for {
			fmt.Print("Enter expression: ")
			expression, _ := reader.ReadString('\n')
			if len(expression) == 0 {
				return
			}
			tt = bxp.NewTruthTable(strings.TrimSpace(expression))
			tt.Compute()
			fmt.Println(tt.ToString())
		}
	}
	tt.Compute()
	if *showTable {
		fmt.Println(tt.ToString())
	}
	if *showTrueRows {
		tri := tt.GetTrueRowIndices()
		fmt.Println(tri)
	}
	if *showModel {
		fmt.Printf("Model:\n%s\n", tt.GetExpression().ToString())
	}
	fmt.Print("Press enter to exit")
	_, _ = reader.ReadString('\n')
}
