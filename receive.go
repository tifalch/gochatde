package main

import (
	_ "github.com/mtib/godolta/deltal"
	"net"
	"time"
)

func receive(ip *net.IPAddr, password string) {
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
