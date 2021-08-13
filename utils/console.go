package utils

import (
	"bufio"
	"os"
	"time"
)

var (
	console = bufio.NewReader(os.Stdin)
)

func Readline() string {
	if str, err := console.ReadString('\n'); err != nil {
		return ""
	} else {
		return str
	}
}

func ReadlineWithTimeout(dur time.Duration, def string) string {
	ch := make(chan string)
	defer close(ch)
	go func() {
		select {
		case <-time.After(dur):
		case ch <- Readline():
		}
	}()
	str := def

	select {
	case str = <-ch:
	case <-time.After(dur):
	}

	return str
}
