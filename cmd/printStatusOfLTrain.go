package main

import (
	"bytes"
	"fmt"
	"github.com/aoeu/mta"
	"github.com/jaytaylor/html2text"
	"log"
	"os"
)

func main() {
	service, err := mta.GetSubwayStatus()
	if err != nil {
		log.Fatal(err)
	}
	line, err := mta.GetLine(service, "L")
	out, err := html2text.FromReader(bytes.NewReader([]byte(line.Text)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting HTML to text: %v", err)
		os.Exit(1)
	}

	fmt.Println(out)
}
