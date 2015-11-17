package main

import (
	"github.com/aoeu/mta"
	"flag"
)

func main() {
	var key string
	flag.StringVar(&key, "key", "aoeu", "Open sesame.")
	flag.Parse()
	mta.GetNextLTrains(key)
}
