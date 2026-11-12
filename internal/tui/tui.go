package tui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gleb/go-calc/internal/eval"
)

var buttons = [][]string{
	{"C", "(", ")", "^"},
	{"7", "8", "9", "/"},
	{"4", "5", "6", "*"},
	{"1", "2", "3", "-"},
	{"0", ".", "%", "+"},
	{"="},
}

// UI holds the calculator state for the terminal interface.
type UI struct {
	Expression  string
	Result      string
	HasError    bool
	SelectedRow int
	SelectedCol int
	Precision   int
}

// New creates a new calculator UI state.
func New(precision int) *UI {
	return &UI{
		Precision: precision,
	}
}

// Render writes the calculator screen to the writer.
func (ui *UI) Render(w io.Writer) {
	var b strings.Builder

	b.WriteString("\x1b[H")
	b.WriteString("  ┌──────────────────────────────────────┐\n")
	b.WriteString("  │            STANDARD CALC             │\n")
	b.WriteString("  ├──────────────────────────────────────┤\n")

	expr := ui.Expression
	if len(expr) > 30 {
		expr = "..." + expr[len(expr)-27:]
	}
	b.WriteString("  │ " + padRight(expr, 36) + " │\n")

	res := ui.Result
	if res == "" {
		res = "0"
	}
	if len(res) > 30 {
		res = res[:27] + "..."
	}

	if ui.HasError {
		b.WriteString("  │ \x1b[1;31m= " + padRight(res, 34) + "\x1b[0m │\n")
	} else {
		b.WriteString("  │ \x1b[1;32m= " + padRight(res, 34) + "\x1b[0m │\n")
	}
	b.WriteString("  ├──────────────────────────────────────┤\n")

	for r := 0; r < len(buttons); r++ {
		b.WriteString("  │ ")
		if r == 5 {
			label := padRight("[        =        ]", 36)
			if ui.SelectedRow == 5 {
				b.WriteString("\x1b[7;1;33m" + label + "\x1b[0m")
			} else {
				b.WriteString("\x1b[1;36m" + label + "\x1b[0m")
			}
		} else {
			for c := 0; c < 4; c++ {
				btn := buttons[r][c]
				label := padRight("[ "+btn+" ]", 7)
				if ui.SelectedRow == r && ui.SelectedCol == c {
					b.WriteString("\x1b[7;1;33m" + label + "\x1b[0m  ")
				} else {
					b.WriteString("\x1b[1;37m" + label + "\x1b[0m  ")
				}
			}
		}
		b.WriteString(" │\n")
	}

	b.WriteString("  └──────────────────────────────────────┘\n")
	b.WriteString("  \x1b[2mType numbers/operators directly on keyboard\x1b[0m\n")
	b.WriteString("  \x1b[2mArrows: Navigate | Space/Enter: Press | q: Quit\x1b[0m\n")

	_, _ = io.WriteString(w, b.String())
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// HandleInput processes a keyboard event.
func (ui *UI) HandleInput(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	if data[0] == 0x1b {
		if len(data) == 1 {
			return true
		}
		if len(data) >= 3 && data[1] == '[' {
			switch data[2] {
			case 'A':
				if ui.SelectedRow > 0 {
					ui.SelectedRow--
				}
			case 'B':
				if ui.SelectedRow < len(buttons)-1 {
					ui.SelectedRow++
				}
			case 'C':
				if ui.SelectedRow < 5 && ui.SelectedCol < 3 {
					ui.SelectedCol++
				}
			case 'D':
				if ui.SelectedCol > 0 {
					ui.SelectedCol--
				}
			}
			return false
		}
		return false
	}

	if data[0] == 3 || data[0] == 'q' || data[0] == 'Q' {
		return true
	}

	if data[0] == '\r' || data[0] == '\n' || data[0] == ' ' {
		ui.pressSelectedButton()
		return false
	}

	if data[0] == 127 || data[0] == 8 {
		if len(ui.Expression) > 0 {
			ui.Expression = ui.Expression[:len(ui.Expression)-1]
			ui.evalLive()
		}
		return false
	}

	if data[0] == 'c' || data[0] == 'C' {
		ui.Expression = ""
		ui.Result = ""
		ui.HasError = false
		return false
	}

	if data[0] == '=' {
		ui.evaluate()
		return false
	}

	ch := data[0]
	if (ch >= '0' && ch <= '9') || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' || ch == '^' || ch == '(' || ch == ')' || ch == '.' {
		ui.Expression += string(ch)
		ui.evalLive()
		return false
	}

	return false
}

func (ui *UI) pressSelectedButton() {
	if ui.SelectedRow == 5 {
		ui.evaluate()
		return
	}

	btn := buttons[ui.SelectedRow][ui.SelectedCol]
	if btn == "C" {
		ui.Expression = ""
		ui.Result = ""
		ui.HasError = false
		return
	}

	ui.Expression += btn
	ui.evalLive()
}

func (ui *UI) evalLive() {
	if strings.TrimSpace(ui.Expression) == "" {
		ui.Result = "0"
		ui.HasError = false
		return
	}

	val, err := eval.Eval(ui.Expression)
	if err == nil {
		ui.Result = eval.Format(val, ui.Precision)
		ui.HasError = false
	}
}

func (ui *UI) evaluate() {
	if strings.TrimSpace(ui.Expression) == "" {
		ui.Result = "0"
		ui.HasError = false
		return
	}

	val, err := eval.Eval(ui.Expression)
	if err != nil {
		ui.Result = err.Error()
		ui.HasError = true
	} else {
		ui.Result = eval.Format(val, ui.Precision)
		ui.HasError = false
		ui.Expression = ui.Result
	}
}

// Run starts the interactive terminal UI.
func Run(precision int) error {
	cmd := exec.Command("stty", "cbreak", "-echo")
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return errors.New("TUI mode requires an interactive terminal (TTY)")
	}

	defer func() {
		restoreCmd := exec.Command("stty", "-cbreak", "echo")
		restoreCmd.Stdin = os.Stdin
		_ = restoreCmd.Run()
		fmt.Print("\x1b[?25h\x1b[0m\n")
	}()

	fmt.Print("\x1b[?25l\x1b[2J")

	ui := New(precision)
	ui.Render(os.Stdout)

	buf := make([]byte, 16)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			continue
		}

		if ui.HandleInput(buf[:n]) {
			break
		}
		ui.Render(os.Stdout)
	}

	return nil
}
