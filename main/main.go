package main

import (
	"fmt"
)

func main() {

	w := NewWriter()
	fmt.Fprint(w, "abc")
	fmt.Fprint(w, "def")

	fmt.Println(string(<-w.channel))
	fmt.Println(string(<-w.channel))
}

type writer struct {
	channel chan []byte
}

func NewWriter() *writer {
	return &writer{channel: make(chan []byte, 100)}
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.channel <- p
	//fmt.Println(&w)
	return len(p), nil
}
