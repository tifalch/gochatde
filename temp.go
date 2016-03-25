package main

// used for temporary testing functions

import (
	"flag"
	"fmt"
	"github.com/mtib/mtiblib/probability"
	"time"
)

func enterPassword() string {
	probability.Initialize()
	setColor(red)
	fmt.Print("password >> ")
	setColor(green)
	for n := 0; n < probability.RandomI(5, 10); n++ {
		fmt.Print("*")
		time.Sleep(time.Duration(50+probability.RandomI(10, 100)) * time.Millisecond)
	}
	resetColor()
	fmt.Print("\n")
	return ""
}

func print(a ...interface{}) {
	if *debug {
		fmt.Print(a...)
	}
}

func printf(f string, a ...interface{}) {
	if *debug {
		fmt.Printf(f, a...)
	}
}

func flagPrint() {
	for k := range flag.Args() {
		printf("flag.Args(%d) == %s\n", k, flag.Arg(k))
	}
	printf("color:    %5v / false\n", *useColor)
	printf("compress: %5v / false\n", *useCompress)
	printf("checksum: %5v / true\n", *useChecksum)
}
