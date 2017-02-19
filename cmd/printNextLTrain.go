package main

import (
	"flag"
	"github.com/aoeu/mta"
	"log"
	"os"
)

func main() {
	var key string
	flag.StringVar(&key, "key", "", "Open sesame.")
	flag.Parse()
	if key == "" {
		flag.Usage()
		os.Exit(1)
	}
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
