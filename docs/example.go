package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/peick/go-gron"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Open("docs/example.json")
	check(err)

	reader := bufio.NewReader(f)

	g := gron.New(reader, gron.OriginalGronFormatter())

	out, err := g.String()
	check(err)

	fmt.Println(out)
}
