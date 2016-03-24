package main

import (
	"fmt"
	"testing"
)

func TestIPs(t *testing.T) {
	fmt.Println("start ip testing")
	fmt.Println(toIP("127.0.0.1"))
	fmt.Println(toIP("127.0.0.1:9999"))
}
