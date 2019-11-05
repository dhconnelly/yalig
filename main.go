package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("== yalig!")
	var b *bufio.Reader
	switch len(os.Args) {
	case 1:
		b = bufio.NewReader(os.Stdin)
	case 2:
		in, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		b = bufio.NewReader(in)
	default:
		log.Fatal("Usage: yalig [file]")
	}
	p := NewParser(NewLexer(b))
	e := NewEvaluator()
	for {
		// Read
		expr, err := p.Parse()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		// Eval
		_, err = e.Eval(expr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
