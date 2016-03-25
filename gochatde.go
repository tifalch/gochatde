package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mtib/godolta/deltal"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
)

var (
	port  = int16(15327)
	hello = `Welcome to gochatde
an encrypted terminal chat client using delta-l encryption.`

	notarget = fmt.Sprintf(`You need to tell gochatde an ip[:port] to connect to
the default port is %d

Usage:
	gochatde [-color -gzip] IP
	gochatde [-color -gzip -check=false] IP:PORT
`, port)
	useColor      = flag.Bool("color", false, "use color")
	useCompress   = flag.Bool("gzip", false, "use compression")
	useChecksum   = flag.Bool("check", true, "use checksum")
	debug         = flag.Bool("debug", false, "run in debug mode")
	chatWaitgroup = new(sync.WaitGroup)
	datasend      func(string) error
	osreturn      = make(chan int)
)

func main() {
	flag.BoolVar(useColor, "colour", false, "use color")
	flag.Parse()
	flagPrint()
	args := flag.Args()
	if len(args) == 1 { // IP:PORT and Password
		if ip, err := toIP(args[0]); err != nil { // could read IP:PORT
			print(ip, "\n")
			fmt.Println(hello)
			// read password
			password := enterPassword() // just for now
			chatWaitgroup.Add(1)        // enter chatmode
			if len(args) == 1 {
				go chatmode(ip, password) // without password
			} else {
				go chatmode(ip, password) // with password
			}
			go receive(ip, password)
			chatWaitgroup.Wait()
			os.Exit(<-osreturn)
		} else {
			fmt.Println(err)
		}
	}
	// Something went wrong
	setColor(red)
	fmt.Println(notarget)
	flag.PrintDefaults()
	resetColor()
	os.Exit(1)
}

func chatmode(ip IP, pass string) {
	csig := make(chan os.Signal, 10)
	cerr := make(chan error, 10)
	signal.Notify(csig, os.Kill, os.Interrupt)
	// make connection
	buf := bufio.NewReader(os.Stdin)
	datasend = func(s string) error {
		return send(s, pass)
	}
chatfor:
	for {
		select {
		case err := <-cerr:
			if err != nil {
				setColor(red)
				fmt.Println(err)
				resetColor()
				switch err.(type) {
				case CommandError:
					if int(err.(CommandError)) == 201 {
						defer func() { osreturn <- 0 }()
					} else {
						defer func() { osreturn <- 2 }()
					}
				default:
					defer func() { osreturn <- 2 }()
				}
				break chatfor
			}
		case <-csig:
			// send "bye"
			defer func() { osreturn <- 0 }()
			break chatfor
		default:
			setColor(green)
			fmt.Print(" Δ ")
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
	chatWaitgroup.Done()
}

func send(message, password string) error {
	msgBuffer := NewMessageBuffer(message)
	var reader *deltal.Encoder
	var err error
	setColor(cyan)
	if *useCompress {
		reader, err = deltal.NewCompressedEncoderReader(msgBuffer, password, *useChecksum)
		print("--- begins compressed enrypted data ---\n")
	} else {
		reader, err = deltal.NewEncoderReader(msgBuffer, password, *useChecksum)
		print("--- begins enrypted data ---\n")
	}
	var data []byte
	b := make([]byte, 12)
	for i := 1; true; i++ {
		n, err2 := reader.Read(b)
		printf("[%02X] (%02d) %s\n", i, n, toString(b, n))
		data = append(data, b[:n]...)
		if err2 != nil {
			break
		}
	}
	print("len(data) == ", len(data), " bytes\n")
	if *useCompress {
		print("--- end of compressed enrypted data ---\n")
	} else {
		print("--- end of enrypted data ---\n")
	}
	resetColor()
	return err
}

func toString(b []byte, read int) string {
	s := "["
	n := 1
	if b[0] == 0xCE && b[1] == 0x94 && b[2] == 0x4C && b[3] == 0xA {
		s += " CHECKSUM:  "
		n = 5
	}
	s += fmt.Sprintf("%02X", b[n-1])
	for k, v := range b[n:] {
		r := ' '
		if n+k == read {
			r = '|'
		}
		s = fmt.Sprintf("%s%c%02X", s, r, v)
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

var commandpairs = map[string]string{
	"§bye, §quit":  "end gochatde",
	"§file <file>": "send file",
	"§ls":          "list files current directory",
	"§pwd":         "prints current working directory",
	"§cd <dir>":    "changes directory",
}

func handleCommand(cmd string) error {
	file := "file"
	setColor(cyan)
	cmd = strings.TrimSpace(cmd)
	switch cmd { // switch all cmds (but handle only single word ones)
	case "bye", "quit":
		return CommandError(201)
	case "help", "?":
		for cmd, expl := range commandpairs {
			fmt.Printf("%s\n\t%s\n", cmd, expl)
		}
		return nil
	case file:
		return handleCommand("help")
	case "ls":
		info, err := ioutil.ReadDir(".")
		if err != nil {
			fmt.Println(err)
			return nil
		}
		var max struct{ Name, Size int }
		for _, f := range info {
			name := len(f.Name())
			size := len(fmt.Sprintf("%d", f.Size))
			if name > max.Name {
				max.Name = name
			}
			if size > max.Size {
				max.Size = size
			}
		}
		format := fmt.Sprintf("%%-%ds %%%dd %%s\n", max.Name, max.Size)
		for _, f := range info {
			if f.IsDir() {
				file = "dir"
			}
			fmt.Printf(format, f.Name(), f.Size(), file)
			file = "file"
		}
		return nil
	case "cd":
		fmt.Println("Usage: §cd <dir>")
		return nil
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(pwd)
		}
		return nil
	} // end single word cmds
	if cmdpart := strings.Fields(cmd); len(cmdpart) > 1 { // begin multi word cmds
		switch cmdpart[0] {
		case file:
			data, err := ioutil.ReadFile(cmdpart[1])
			if err != nil {
				fmt.Println(err)
				return nil
			}
			return datasend(string(data))
		case "cd":
			err := os.Chdir(cmdpart[1])
			if err != nil {
				fmt.Println("directory", cmdpart[1], "was not found.")
			} else {
				fmt.Println("switched directory")
			}
			return nil
		}
	} // end multi word cmd
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
