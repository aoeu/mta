package main

import (
	"flag"
	"github.com/aoeu/mta"
	"log"
	"time"
)

func main() {
	var key string
	flag.StringVar(&key, "key", "aoeu", "Open sesame.")
	flag.Parse()
	start := 8
	stop := 11
	updating := false
	quit := make(chan bool)
	for {
		h := time.Now().Hour()
		switch {
		case !updating && h > start && h < stop:
			updating = true
			go getTrainTimes(key, quit)
		case updating && h >= stop:
			updating = false
			quit <- true
		}
		<-time.After(1 * time.Minute)
	}
}

func getTrainTimes(key string, quit <-chan bool) {
	log.Println("Starting to get train times.")
	defer log.Println("Done getting train times.")
	go getAndLogTrainTime(key)
	for {
		select {
		case <-time.After(10 * time.Second):
			go getAndLogTrainTime(key)
		case <-quit:
			return
		}
	}

}

func getAndLogTrainTime(key string) {
	t, err := mta.GetNextMontroseLTrains(key)
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range t[0:1] {
		log.Println(s)
	}

}
