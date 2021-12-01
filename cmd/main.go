package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/guest42069/i6r9c/connection"
	"github.com/guest42069/i6r9c/msg"
	"github.com/guest42069/i6r9c/worker"
	"golang.org/x/term"
)

// PrintMsg will print out a formatted Msg m to the provided Terminal t.
func PrintMsg(t *term.Terminal, m *msg.Msg) {
	var line string
	switch m.Cmd() {
	case "PRIVMSG":
		line = fmt.Sprintf("[%s] [%s%s%s@%s%s%s] %s",
			m.Timestamp(),
			string(t.Escape.Yellow), m.Nick(), string(t.Escape.Reset),
			string(t.Escape.Yellow), m.Rcpt(), string(t.Escape.Reset),
			m.Content(),
		)
	case "JOIN":
		line = fmt.Sprintf("[%s] [%s!%s@%s] has joined [%s%s]",
			m.Timestamp(),
			m.Nick(),
			m.User(),
			m.Host(),
			m.Rcpt(),
			m.Content(),
		)
	case "PART":
		line = fmt.Sprintf("[%s] [%s!%s@%s] has parted [%s: %s]",
			m.Timestamp(),
			m.Nick(),
			m.User(),
			m.Host(),
			m.Rcpt(),
			m.Content(),
		)
	case "QUIT":
		line = fmt.Sprintf("[%s] [%s!%s@%s] has quit [%s]",
			m.Timestamp(),
			m.Nick(),
			m.User(),
			m.Host(),
			m.Content(),
		)
	case "NOTICE":
		line = fmt.Sprintf("[%s] [%s%s%s@%s%s%s] %s%s%s [%s%s%s]",
			m.Timestamp(),
			string(t.Escape.Yellow), m.Nick(), string(t.Escape.Reset),
			string(t.Escape.Yellow), m.Rcpt(), string(t.Escape.Reset),
			string(t.Escape.Yellow), m.Content(), string(t.Escape.Reset),
			string(t.Escape.Yellow), m.Args(), string(t.Escape.Reset),
		)
	case "ERROR":
		line = fmt.Sprintf("[%s] [%s%s%s@%s%s%s] %s%s%s [%s%s%s]",
			m.Timestamp(),
			string(t.Escape.Red), m.Nick(), string(t.Escape.Reset),
			string(t.Escape.Red), m.Rcpt(), string(t.Escape.Reset),
			string(t.Escape.Red), m.Content(), string(t.Escape.Reset),
			string(t.Escape.Red), m.Args(), string(t.Escape.Reset),
		)
	case "NICK":
		line = fmt.Sprintf("[%s] [%s%s%s] changed renamed to [%s%s%s]",
			m.Timestamp(),
			string(t.Escape.Yellow), m.Nick(), string(t.Escape.Reset),
			string(t.Escape.Red), m.Content(), string(t.Escape.Reset),
		)
	default:
		line = fmt.Sprintf("[%s] [%s%s%s] [%s%s%s@%s%s%s] %s%s%s [%s%s%s]",
			m.Timestamp(),
			string(t.Escape.Yellow), m.Cmd(), string(t.Escape.Reset),
			string(t.Escape.Yellow), m.Nick(), string(t.Escape.Reset),
			string(t.Escape.Yellow), m.Rcpt(), string(t.Escape.Reset),
			string(t.Escape.Yellow), m.Content(), string(t.Escape.Reset),
			string(t.Escape.Yellow), m.Args(), string(t.Escape.Reset),
		)
	}
	t.Write([]byte(line + "\n"))
}

// setPrompt will update the Terminal prompt with the provided value.
func setPrompt(t *term.Terminal, prompt string) {
	t.SetPrompt(fmt.Sprintf("%s%s%s> ",
		string(t.Escape.Yellow), prompt, string(t.Escape.Reset),
	))
}

func main() {
	serverPtr := flag.String("server", "ircs://irc.oftc.net:6697/", "URL schema of server, [scheme]://[server]:[port]. irc for non-TLS, ircs for TLS.")
	proxyPtr := flag.String("proxy", "socks5://127.0.0.1:9050/", "URL schema of proxy, [scheme]://[server]:[port].")
	nickPtr := flag.String("nick", "", "IRC nickname to use.")
	saslPtr := flag.String("sasl", "", "SASL cert and key prefix (I.E foo/bar for foo/bar.crt and foo/bar.key)")
	flag.Parse()
	if len(*nickPtr) < 1 {
		nickPtr = randomName()
	}
	var conn net.Conn
	var err error
	if len(*saslPtr) > 0 {
		cert, err := tls.LoadX509KeyPair(*saslPtr+".crt", *saslPtr+".key")
		if err != nil {
			panic(err)
		}
		conn, err = connection.Connect(proxyPtr, serverPtr, &cert, true)
		if err != nil {
			panic(err)
		}
	} else {
		conn, err = connection.Connect(proxyPtr, serverPtr, nil, true)
		if err != nil {
			panic(err)
		}

	}
	recv, send, stop := worker.Worker(conn)
	connection.Login(send, *nickPtr, (len(*saslPtr) > 0))
	defer close(stop)
	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		panic("not a terminal")
	}
	old, err := term.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer term.Restore(0, old)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	t := term.NewTerminal(screen, "> ")
	go func() {
		for {
			select {
			case <-stop:
				return
			case m := <-recv:
				if m.Cmd() == "PING" {
					send <- fmt.Sprintf("PONG :%s", m.Content())
				} else {
					PrintMsg(t, m)
				}
			}
		}
	}()
	var current string = ""
	for {
		line, err := t.ReadLine()
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(line, "/") {
			line = line[1:]
			arr := strings.SplitN(line, " ", 2)
			if len(arr) > 1 {
				line = arr[1]
			} else {
				line = ""
			}
			switch arr[0] {
			case "msg":
				arr = strings.SplitN(line, " ", 2)
				current = arr[0]
				setPrompt(t, current)
				if len(arr) > 1 {
					send <- fmt.Sprintf("PRIVMSG %s :%s", arr[0], arr[1])
				}
			case "join":
				arr = strings.SplitN(line, " ", 2)
				current = arr[0]
				setPrompt(t, current)
				send <- fmt.Sprintf("JOIN %s", arr[0])
			case "part":
				arr = strings.SplitN(line, " ", 2)
				send <- fmt.Sprintf("PART %s", arr[0])
			case "quit":
				send <- fmt.Sprintf("QUIT :%s", line)
				return
			case "nick":
				send <- fmt.Sprintf("NICK %s", line)
			}
		} else {
			send <- fmt.Sprintf("PRIVMSG %s :%s", current, line)
		}
	}
}
