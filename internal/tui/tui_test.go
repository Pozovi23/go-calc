package tui_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gleb/go-calc/internal/tui"
)

func TestUI_TypingAndEvaluation(t *testing.T) {
	ui := tui.New(10)

	// Type "2+3"
	ui.HandleInput([]byte{'2'})
	ui.HandleInput([]byte{'+'})
	ui.HandleInput([]byte{'3'})

	if ui.Expression != "2+3" {
		t.Fatalf("expected expression '2+3', got %q", ui.Expression)
	}

	// Press '='
	ui.HandleInput([]byte{'='})
	if ui.Result != "5" {
		t.Errorf("expected result '5', got %q", ui.Result)
	}
}

func TestUI_ArrowNavigationAndButtonPress(t *testing.T) {
	ui := tui.New(10)

	// Move down to row 1 (button '7')
	ui.HandleInput([]byte{0x1b, '[', 'B'}) // Down
	if ui.SelectedRow != 1 {
		t.Fatalf("expected selected row 1, got %d", ui.SelectedRow)
	}

	// Press Space to click '7'
	ui.HandleInput([]byte{' '})
	if ui.Expression != "7" {
		t.Errorf("expected expression '7', got %q", ui.Expression)
	}

	// Move right to button '8'
	ui.HandleInput([]byte{0x1b, '[', 'C'}) // Right
	if ui.SelectedCol != 1 {
		t.Fatalf("expected selected col 1, got %d", ui.SelectedCol)
	}

	// Press Enter to click '8'
	ui.HandleInput([]byte{'\n'})
	if ui.Expression != "78" {
		t.Errorf("expected expression '78', got %q", ui.Expression)
	}
}

func TestUI_BackspaceAndClear(t *testing.T) {
	ui := tui.New(10)

	// Type "123"
	ui.HandleInput([]byte{'1'})
	ui.HandleInput([]byte{'2'})
	ui.HandleInput([]byte{'3'})

	// Backspace
	ui.HandleInput([]byte{127})
	if ui.Expression != "12" {
		t.Errorf("expected '12' after backspace, got %q", ui.Expression)
	}

	// Clear 'c'
	ui.HandleInput([]byte{'c'})
	if ui.Expression != "" {
		t.Errorf("expected empty expression after clear, got %q", ui.Expression)
	}
}

func TestUI_Exit(t *testing.T) {
	ui := tui.New(10)

	if !ui.HandleInput([]byte{'q'}) {
		t.Error("expected exit true on 'q'")
	}
	if !ui.HandleInput([]byte{0x1b}) { // Esc
		t.Error("expected exit true on bare Esc")
	}
	if !ui.HandleInput([]byte{3}) { // Ctrl+C
		t.Error("expected exit true on Ctrl+C")
	}
}

func TestUI_Render(t *testing.T) {
	ui := tui.New(10)
	ui.Expression = "2+2"
	ui.Result = "4"

	var buf bytes.Buffer
	ui.Render(&buf)

	out := buf.String()
	if !strings.Contains(out, "STANDARD CALC") {
		t.Error("render output does not contain header")
	}
	if !strings.Contains(out, "2+2") {
		t.Error("render output does not contain expression")
	}
}
