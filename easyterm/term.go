package easyterm

import (
	"fmt"
)

var (
	// CSI is escape character, ESC
	CSI = '\x1b'
)

// Rewrite returns sequence that erases the current line, writes new string
func Rewrite(new string) string {
	return fmt.Sprintf("\x1b[G%s\x1b[K", new)
}

// ClearScreen returns sequence that clears the terminal, puts cursor on 1:1
func ClearScreen() string {
	return fmt.Sprint("\x1b[J")
}

// Position returns sequence that puts cursor on x:y
func Position(x, y uint) string {
	return fmt.Sprintf("\x1b[%d;%df", x, y)
}

// Italics returns sequence that will write following characters in italics
func Italics() string {
	return "\x1b[3m"
}

// Reset returns sequence that resets settings
func Reset() string {
	return "\x1b[0m"
}

// Bold returns sequence that turns text bold
func Bold() string {
	return "\x1b[1m"
}

// Underline control sequence
func Underline() string {
	return "\x1b[4m"
}

// Negative control sequence
func Negative() string {
	return "\x1b[7m"
}

// Strike control sequence
func Strike() string {
	return "\x1b[9m"
}

// Style contains all format information needed to print with style!
type Style struct {
	Bold, Italics, Underline, Negative, Strike bool
	Foreground, Background                     Color
	ResetAfter                                 bool
}

func (sty *Style) Write(s string) string {
	ret := ""
	if sty.Bold {
		ret += Bold()
	}
	if sty.Italics {
		ret += Italics()
	}
	if sty.Underline {
		ret += Underline()
	}
	if sty.Negative {
		ret += Negative()
	}
	if sty.Strike {
		ret += Strike()
	}
	if sty.Foreground != nil {
		ret += sty.Foreground.Control(Foreground)
	}
	if sty.Background != nil {
		ret += sty.Background.Control(Background)
	}
	ret += s
	if sty.ResetAfter {
		ret += Reset()
	}
	return ret
}

// NewStyle creates default style
func NewStyle() *Style {
	return new(Style)
}
