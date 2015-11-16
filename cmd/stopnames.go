package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/aoeu/mta"
	"io"
	"log"
	"os"
	"regexp"
)

var re = regexp.MustCompile("^L.*")

func main() {
	in := "/home/tasm/ir/src/github.com/aoeu/mta/data/static/stops.txt"
	out := "/home/tasm/ir/src/github.com/aoeu/mta/data/static/stops_l_train.json"
	f, err := os.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f)
	stops := make(mta.StopNames, 0)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if re.MatchString(record[0]) {
			stops[mta.StopID(record[0])] = mta.StopName(record[2])
		}
	}
	b, err := json.Marshal(stops)
	if err != nil {
		log.Fatal(err)
	}
	f, err = os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
