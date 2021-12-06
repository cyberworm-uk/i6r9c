package worker

import (
	"bufio"
	"net"
	"strings"
	"sync"

	"github.com/guest42069/i6r9c/msg"
)

// Worker will return 3 channels: output, input and a stop channel. Messages from the IRCd on the other side of conn are returned via the output channel, lines to send to the remote server should be sent to the input channel. closing the stop channel will stop the workers.
func Worker(conn net.Conn, wg *sync.WaitGroup) (<-chan *msg.Msg, chan<- string, chan bool) {
	toCaller := make(chan *msg.Msg)
	fromCaller := make(chan string)
	lines := make(chan string)
	stopWorker := make(chan bool)
	wg.Add(2)
	go func(rd *bufio.Reader) {
		for {
			line, err := rd.ReadString('\n')
			if err != nil {
				close(stopWorker)
				return
			}
			line = strings.TrimRight(line, " \r\n")
			lines <- line
		}
	}(bufio.NewReader(conn))
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopWorker:
				return
			case line := <-lines:
				toCaller <- msg.Parse(line)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopWorker:
				return
			case line := <-fromCaller:
				_, err := conn.Write([]byte(line + "\n"))
				if err != nil {
					close(stopWorker)
					return
				}
			}
		}
	}()
	return toCaller, fromCaller, stopWorker
}
