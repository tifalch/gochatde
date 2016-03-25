package easyterm

import (
	"fmt"
)

type Color interface {
	Control(foreground bool) string
}

type RGB struct {
	R, G, B byte
}

type ANSI uint8

func (c *RGB) Control(fg bool) string {
	num := 38
	if !fg {
		num = 48
	}
	return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", num, c.R, c.G, c.B)
}

func (c *ANSI) Control(fg bool) string {
	num := int(*c) + 30
	if !fg {
		num += 10
	}
	return fmt.Sprintf("\x1b[%dm", num)
}

var (
	Black   ANSI
	Red     ANSI = 1
	Green   ANSI = 2
	Yellow  ANSI = 3
	Blue    ANSI = 4
	Magenta ANSI = 5
	Cyan    ANSI = 6
	White   ANSI = 7
)

const (
	Background = false
	Foreground = true
)

func NewRGB(r, g, b int) *RGB {
	return &RGB{byte(r), byte(g), byte(b)}
}
