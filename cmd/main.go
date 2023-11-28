package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/cyberworm-uk/i6r9c/connection"
	"github.com/cyberworm-uk/i6r9c/msg"
	"github.com/cyberworm-uk/i6r9c/worker"
	"golang.org/x/term"
)

func updateTerminalSize(t *term.Terminal) {
	if w, h, err := term.GetSize(int(os.Stdin.Fd())); err != nil {
		t.SetSize(w, h)
	}
}

// PrintMsg will print out a formatted Msg m to the provided Terminal t.
func printMsg(t *term.Terminal, m *msg.Msg) {
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
			string(t.Escape.Yellow), m.Content(), string(t.Escape.Reset),
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

func printUsage(t *term.Terminal) {
	t.Write([]byte(fmt.Sprintf("%s%s%s %s%s%s %s\n",
		string(t.Escape.White), "/msg", string(t.Escape.Reset),
		string(t.Escape.White), "<#channel/recipient>", string(t.Escape.Reset),
		"[message]",
	)))
	t.Write([]byte(fmt.Sprintf("%s%s%s %s%s%s\n",
		string(t.Escape.White), "/join", string(t.Escape.Reset),
		string(t.Escape.White), "<#channel>", string(t.Escape.Reset),
	)))
	t.Write([]byte(fmt.Sprintf("%s%s%s %s%s%s %s\n",
		string(t.Escape.White), "/part", string(t.Escape.Reset),
		string(t.Escape.White), "<#channel>", string(t.Escape.Reset),
		"[reason]",
	)))
	t.Write([]byte(fmt.Sprintf("%s%s%s %s%s%s\n",
		string(t.Escape.White), "/nick", string(t.Escape.Reset),
		string(t.Escape.White), "<newnick>", string(t.Escape.Reset),
	)))
	t.Write([]byte(fmt.Sprintf("%s%s%s %s\n",
		string(t.Escape.White), "/quit", string(t.Escape.Reset),
		"[reason]",
	)))
}

func sendMsg(current, content string, send chan<- string) {
	var frag string
	for {
		if len(content) > 0 {
			if len(content) > 400 {
				frag = content[:400]
				content = content[400:]
			} else {
				frag = content
				content = ""
			}
			send <- fmt.Sprintf("PRIVMSG %s :%s", current, frag)
		} else {
			return
		}
	}
}

func main() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	serverPtr := flag.String("server", "ircs://irc.oftc.net:6697/", "URL schema of server, [scheme]://[server]:[port]. irc for non-TLS, ircs for TLS.")
	proxyPtr := flag.String("proxy", "socks5://127.0.0.1:9050/", "URL schema of proxy, [scheme]://[server]:[port].")
	nickPtr := flag.String("nick", "", "IRC nickname to use.")
	saslPtr := flag.String("sasl", "", "SASL cert and key prefix (I.E foo/bar for foo/bar.crt and foo/bar.key)")
	verifyPtr := flag.Bool("verify", true, "Verify TLS certificates (I.E. an .onion with TLS but no valid cert.)")
	flag.Parse()
	if len(*nickPtr) < 1 {
		nickPtr = randomName()
	}
	var conn net.Conn
	var err error
	if len(*saslPtr) > 0 {
		cert, err := tls.LoadX509KeyPair(*saslPtr+".crt", *saslPtr+".key")
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}
		conn, err = connection.Connect(proxyPtr, serverPtr, &cert, *verifyPtr)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}
	} else {
		conn, err = connection.Connect(proxyPtr, serverPtr, nil, *verifyPtr)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

	}
	recv, send, stop := worker.Worker(conn, wg)
	connection.Login(send, *nickPtr, (len(*saslPtr) > 0))
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf("ERROR: not a terminal\n")
		return
	}
	old, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), old)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	t := term.NewTerminal(screen, "> ")
	resizeChan := make(chan os.Signal)
	go func() {
		for range resizeChan {
			updateTerminalSize(t)
		}
	}()
	signal.Notify(resizeChan, syscall.SIGWINCH)
	go func() {
		for {
			select {
			case <-stop:
				return
			case m := <-recv:
				if m.Cmd() == "PING" {
					send <- fmt.Sprintf("PONG :%s", m.Content())
				} else {
					printMsg(t, m)
				}
			}
		}
	}()
	lineChan := make(chan string)
	go func() {
		for {
			if line, err := t.ReadLine(); err != nil {
				printMsg(t, msg.InternalError(err))
				close(stop)
				return
			} else {
				lineChan <- line
			}
		}
	}()
	var current string = ""
	for {
		select {
		case <-stop:
			return
		case line := <-lineChan:
			if strings.HasPrefix(line, "/") {
				line = line[1:]
				cmd, line := msg.Split(line, " ")
				switch cmd {
				case "msg":
					rcpt, content := msg.Split(line, " ")
					current = rcpt
					setPrompt(t, current)
					if len(content) > 0 {
						sendMsg(current, content, send)
					}
				case "join":
					channel, _ := msg.Split(line, " ")
					current = channel
					setPrompt(t, current)
					send <- fmt.Sprintf("JOIN %s", channel)
				case "part":
					channel, content := msg.Split(line, " ")
					send <- fmt.Sprintf("PART %s :%s", channel, content)
				case "quit":
					send <- fmt.Sprintf("QUIT :%s", line)
					close(stop)
					return
				case "nick":
					send <- fmt.Sprintf("NICK %s", line)
				default:
					printUsage(t)
				}
			} else {
				sendMsg(current, line, send)
			}
		}
	}
}
