package easyterm

import (
	"fmt"
	"testing"
)

func TestColors(t *testing.T) {
	ansi := []ANSI{Black, Red, Green, Yellow, Blue, Magenta, Cyan, White}
	ansiname := []string{"Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White"}
	for k := range ansi {
		fmt.Println(ansi[k].Control(Background), ansiname[k], Reset(), ansi[k].Control(Foreground), ansiname[k], Reset())
	}
	for k := 0; k < 512; k += 8 {
		if k%128 == 0 {
			fmt.Print("\n")
		}
		fmt.Print(NewRGB(k, k<<2, k<<4).Control(true), "Color", Reset())
	}
	fmt.Print("\n")
}
