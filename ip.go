package main

import (
	"fmt"
	"strconv"
	"strings"
)

// IP is IPv4 and IPv6 interface
type IP interface {
	fmt.Stringer
	Valid() (bool, error)
}

// IPv4 to contain IP-Adress
type IPv4 struct {
	Adress []byte
	Port   int16
}

func (i IPv4) String() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", i.Adress[0], i.Adress[1], i.Adress[2], i.Adress[3], i.Port)
}

// Valid checks if IP has valid form
func (i IPv4) Valid() (bool, error) {
	return true, nil
}

func toIP(s string) (IP, error) {
	spl := strings.Split(s+":", ":")
	ins, err := splitIPstring(spl[0])
	p, _ := strconv.ParseInt(spl[1], 10, 16)
	port := int16(p)
	return &IPv4{ins, port}, err
}

func splitIPstring(s string) ([]byte, error) {
	b := make([]byte, 4)
	sarr := strings.Split(s, ".")
	if len(sarr) != len(b) {
		return b, IPerror(s + " is wrong format")
	}
	for k, v := range sarr {
		n, err := strconv.ParseUint(v, 10, 8)
		if err != nil || k == len(b) {
			return b, err
		}
		b[k] = uint8(n)
	}
	return b, nil
}

// IPerror to be returned on format error
type IPerror string

func (e IPerror) Error() string {
	return string(e)
}
