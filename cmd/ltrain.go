package main

import (
	"flag"
	"github.com/aoeu/mta"
	"log"
)

func main() {
	var key string
	flag.StringVar(&key, "key", "aoeu", "Open sesame.")
	flag.Parse()
	t, err := mta.GetNextMontroseLTrains(key)
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range t {
		log.Println(s)
	}
}
