package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
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
	for {
		expr, err := p.Parse()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		fmt.Println(expr)
	}
}
