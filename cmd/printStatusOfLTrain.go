package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/aoeu/mta"
	"github.com/jaytaylor/html2text"
	"log"
	"os"
	"strings"
)

func main() {
	service, err := mta.GetSubwayStatus()
	if err != nil {
		log.Fatal(err)
	}
	line, err := mta.GetLine(service, "L")
	status, err := html2text.FromReader(bytes.NewReader([]byte(line.Text)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting HTML to text: %v", err)
		os.Exit(1)
	}
	switch line.Status {
	case "GOOD SERVICE":
		fmt.Println("Good service into Manhattan")
	case "DELAYS":
		fmt.Println(findTextWith("delays", status, false))
	case "SERVICE CHANGE":
		fmt.Println(findTextWith("trains", status, false))
	case "PLANNED WORK":
		fmt.Println(findTextWith("Planned Work", status, true))
	default:
		fmt.Println(line.Status)
	}
}

func findTextWith(s, in string, lineAfter bool) string {
	sc := bufio.NewScanner(strings.NewReader(in))
	for sc.Scan() {
		if t := sc.Text(); strings.Contains(t, s) {
			if lineAfter {
				sc.Scan()
				return sc.Text()
			}
			return t
		}
	}
	return ""
}
