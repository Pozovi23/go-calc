package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gleb/go-calc/internal/eval"
	"github.com/gleb/go-calc/internal/tui"
)

func main() {
	var precision int
	var showUI bool

	flag.IntVar(&precision, "p", 10, "precision (number of decimal digits)")
	flag.BoolVar(&showUI, "ui", false, "launch visual terminal calculator")
	flag.Parse()

	if precision < 0 {
		fmt.Fprintln(os.Stderr, "error: precision cannot be negative")
		os.Exit(1)
	}

	if showUI {
		if err := tui.Run(precision); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		return
	}

	if flag.NArg() > 0 {
		expr := strings.Join(flag.Args(), " ")
		val, err := eval.Eval(expr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println(eval.Format(val, precision))
		return
	}

	runREPL(precision)
}
