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
	getAndLogTrainTime(key)
}

func getAndLogTrainTime(key string) {
	t, err := mta.GetNextMontroseLTrains(key)
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range t[0:3] {
		log.Println(s)
	}

}
