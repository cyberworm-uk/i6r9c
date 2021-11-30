package worker

import (
	"bufio"
	"net"
	"strings"

	"github.com/guest42069/i6r9c/msg"
)

func Worker(conn net.Conn) (<-chan *msg.Msg, chan<- string, chan bool) {
	toCaller := make(chan *msg.Msg)
	fromCaller := make(chan string)
	stopWorker := make(chan bool)
	go func(rd *bufio.Reader) {
		for {
			select {
			case <-stopWorker:
				return
			default:
				line, err := rd.ReadString('\n')
				if err != nil {
					panic(err)
				}
				line = strings.TrimRight(line, " \r\n")
				toCaller <- msg.Parse(line)
			}
		}
	}(bufio.NewReader(conn))
	go func() {
		for {
			select {
			case line := <-fromCaller:
				_, err := conn.Write([]byte(line + "\n"))
				if err != nil {
					panic(err)
				}
			case <-stopWorker:
				return
			}
		}
	}()
	return toCaller, fromCaller, stopWorker
}
