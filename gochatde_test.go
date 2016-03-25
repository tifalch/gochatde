package main

import (
	"fmt"
	"testing"
)

func TestIPs(t *testing.T) {
	for _, v := range []string{"127.0.0.1", "127.0.0.1:9999"} {
		_, err := toIP(v)
		if err != nil {
			t.Fail()
			fmt.Println("failed at", v)
		}
	}
	for _, v := range []string{"asdb", "192", "....:", "......", ":", ".."} {
		_, err := toIP(v)
		if err == nil {
			t.Fail()
			fmt.Println("validated", v, "... should have failed")
		}
	}
}

// Test if Valid
