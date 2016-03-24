package main

import (
	"flag"
	"fmt"
)

var (
	port  = 15327
	hello = `Welcome to gochatde
an encrypted terminal chat client using delta-l encryption.`

	notarget = fmt.Sprintf("%s %d", `You need to tell gochatde an ip[:port] to connect to
the default port is`, port)
)

func main() {
	flag.Parse()
	args := flag.Args()
	failed := !(len(args) == 1)
	if !failed {
		ip, err := toIP(args[0])
		if err != nil {
			failed = true
			fmt.Println(err)
			goto fail
		}
		b, e := ip.Valid()
		if b == true && e == nil {
			// chatmode
			fmt.Println(hello)
		} else {
			failed = true
		}
	}
fail:
	if failed {
		fmt.Println(notarget)
	}
}
