package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mtib/godolta/deltal"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
)

var (
	port  = 15327
	hello = `Welcome to gochatde
an encrypted terminal chat client using delta-l encryption.`

	notarget = fmt.Sprintf("%s %d", `You need to tell gochatde an ip[:port] to connect to
the default port is`, port)
	useColor = flag.Bool("color", false, "use color")
)

func main() {
	flag.Parse()
	fmt.Println(*useColor)
	args := flag.Args()
	failed := !(len(args) >= 1)
	wg := new(sync.WaitGroup)
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
			wg.Add(1)
			if len(args) == 1 {
				go chatmode(wg, "")
			} else {
				go chatmode(wg, args[1])
			}
		} else {
			failed = true
		}
	}
fail:
	if failed {
		fmt.Println(notarget)
	} else {
		wg.Wait()
	}
}

func chatmode(wg *sync.WaitGroup, pass string) {
	csig := make(chan os.Signal, 10)
	cerr := make(chan error, 10)
	signal.Notify(csig, os.Kill, os.Interrupt)
	// make connection
	buf := bufio.NewReader(os.Stdin)
chatfor:
	for {
		select {
		case err := <-cerr:
			if err != nil {
				setColor(red)
				fmt.Println(err)
				resetColor()
				break chatfor
			}
		case <-csig:
			// send "bye"
			break chatfor
		default:
			setColor(green)
			fmt.Print(" > ")
			var msg []byte
			msg, more, err := buf.ReadLine()
			resetColor()
			if msg == nil || len(msg) == 0 {
				continue
			}
			for more {
				if err != nil {
					panic(err)
				}
				var add []byte
				add, more, err = buf.ReadLine()
				msg = append(msg, add...)
			}
			if msg[0] == 0xC2 && msg[1] == 0xA7 {
				err3 := handleCommand(string(msg[2:]))
				cerr <- err3
				continue // next input
			}
			// else send text
			err2 := send(string(msg), pass)
			cerr <- err
			cerr <- err2
		}
	}
	wg.Done()
}

func send(message, password string) error {
	msgBuffer := NewMessageBuffer(message)
	reader, err := deltal.NewEncoderReader(msgBuffer, password, true)
	setColor(cyan)
	fmt.Println("--- begins encrypted data ---")
	var data []byte
	b := make([]byte, 12)
	for i := 1; true; i++ {
		n, err2 := reader.Read(b)
		fmt.Printf("[%02X] (%02d) %s\n", i, n, toString(b))
		data = append(data, b...)
		if err2 != nil {
			break
		}
	}
	fmt.Println("--- end of encrypted data ---")
	resetColor()
	return err
}

func toString(b []byte) string {
	s := "["
	n := 1
	if b[0] == 0xCE && b[1] == 0x94 && b[2] == 0x4C && b[3] == 0xA {
		s += " CHECKSUM:  "
		n = 5
	}
	s += fmt.Sprintf("%02X", b[n-1])
	for _, v := range b[n:] {
		s = fmt.Sprintf("%s %02X", s, v)
	}
	return s + "]"
}

// CommandError to be thrown when Cmd not working
type CommandError int

func (c CommandError) Error() string {
	switch int(c) {
	case 201:
		return "Bye!"
	case -1:
		fallthrough
	default:
		return "unknown command error"
	}
}

func handleCommand(cmd string) error {
	setColor(cyan)
	cmd = strings.TrimSpace(cmd)
	switch {
	case cmd == "bye", cmd == "quit":
		return CommandError(201)
	}
	setColor(red)
	fmt.Println("unknown command")
	resetColor()
	return nil
}

func resetColor() {
	if *useColor {
		io.WriteString(os.Stdout, "\x1b[0m")
	}
}

func setColor(c color) {
	if *useColor {
		io.WriteString(os.Stdout, fmt.Sprintf("\x1b[%dm", 30+int(c)))
	}
}

type color uint8

const (
	// black   color = iota
	red   color = iota + 1
	green color = iota + 1
	// yellow  color = iota
	// blue color = iota // TODO use this one for output
	// magenta color = iota
	cyan color = iota + 4
	// white   color = iota
)
