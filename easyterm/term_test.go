package easyterm

import (
	"fmt"
	"testing"
)

func TestStyle(t *testing.T) {
	s0 := NewStyle()
	s0.Foreground = NewRGB(255, 0, 0)
	s0.Bold = true
	s0.Underline = true
	s0.ResetAfter = true
	fmt.Print(s0.Write("Hallo Welt\n"))
	fmt.Print("Hallo Welt\n")

	s1 := NewStyle()
	s1.Background = &Blue
	s1.Italics = true
	s1.Negative = true
	s1.Strike = true
	s1.ResetAfter = false
	fmt.Print(s0.Write("Hallo Welt v2\nTest"))
	fmt.Print("Hallo Welt\n")

	fmt.Print(Reset())

	fmt.Print("Wrong")
	fmt.Print(Rewrite("Right!\n"))
}
