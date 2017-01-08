package main

import (
	"flag"
"fmt"
	"github.com/aoeu/mta"
	"log"
	"errors"
)

func main() {
	var key string
	flag.StringVar(&key, "key", "aoeu", "Open sesame.")
	flag.Parse()
	s, err := getNextTrainTimes(key, 3)
        if err != nil {
		log.Fatal(err)
	}
fmt.Print(s)
}

func getNextTrainTimes(key string, quantity int) (string, error) {
	t, err := mta.GetNextMontroseLTrains(key)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get next L train times from feed: %v", err))
	}
        s := ""
	for _, st := range t[0:quantity] {
		s += fmt.Sprintf("%v\n", st.DeltaInUnderTwentyRunes())
	}
	return s, nil
}
