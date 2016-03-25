package main

import (
	_ "github.com/mtib/godolta/deltal"
	_ "net"
	"time"
)

func receive(ip IP, password string) {
	// printf("start receiving from %v", ip)
	for {
		select {
		default:
			time.Sleep(time.Duration(500) * time.Millisecond)
			continue
		}
	}
	panic("receive should never end before main...")
}
