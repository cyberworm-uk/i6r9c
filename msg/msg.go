package msg

import (
	"strings"
	"time"
)

type Msg struct {
	timestamp, nick, user, host, cmd, rcpt, content, args string
}

// User returns the user who sent the message
func (m *Msg) User() string {
	return m.user
}

// Timstamp returns the timestamp the message was parsed at
func (m *Msg) Timestamp() string {
	return m.timestamp
}

// Nick returns the nick who sent the message
func (m *Msg) Nick() string {
	return m.nick
}

// Host returns the host of the user who sent the message.
func (m *Msg) Host() string {
	return m.host
}

// Cmd returns the command associated with the message
func (m *Msg) Cmd() string {
	return m.cmd
}

// Rcpt returns the intended recipient of the message
func (m *Msg) Rcpt() string {
	return m.rcpt
}

// Content returns the message content
func (m *Msg) Content() string {
	return m.content
}

// Args returns any additional arguments associated with the message
func (m *Msg) Args() string {
	return m.args
}

func split(s, d string) (string, string) {
	arr := strings.SplitN(s, d, 2)
	if len(arr) == 2 {
		return arr[0], arr[1]
	} else if len(arr) == 1 {
		return arr[0], ""
	} else {
		return "", ""
	}
}

func Parse(line string) *Msg {
	m := &Msg{}
	m.timestamp = time.Now().Format("15:04:05")
	if strings.HasPrefix(line, ":") {
		line = line[1:]
		m.host, line = split(line, " ")
		m.nick, m.host = split(m.host, "!")
		m.user, m.host = split(m.host, "@")
	}
	line, m.content = split(line, " :")
	m.cmd, line = split(line, " ")
	m.rcpt, m.args = split(line, " ")
	return m
}
