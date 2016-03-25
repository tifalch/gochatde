package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestBuffer(t *testing.T) {
	test := []string{"hallo", "", "\a\t", "¹²³€½¾€³{}"}
	for _, v := range test {
		buff := NewMessageBuffer(v)
		data, err := ioutil.ReadAll(buff)
		if string(data) != v || err != nil {
			t.Fail()
			fmt.Println(v, "!=", string(data))
		} else {
			fmt.Println(v, "==", string(data))
		}
	}
}
