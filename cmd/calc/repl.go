package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gleb/go-calc/internal/eval"
	"github.com/gleb/go-calc/internal/tui"
)

func runREPL(precision int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" || line == "q" {
			break
		}

		if line == "ui" {
			if err := tui.Run(precision); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
			}
			fmt.Print("> ")
			continue
		}

		if line != "" {
			val, err := eval.Eval(line)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
			} else {
				fmt.Println(eval.Format(val, precision))
			}
		}

		fmt.Print("> ")
	}
}
