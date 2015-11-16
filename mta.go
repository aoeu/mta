package mta

import (
	"encoding/json"
	"fmt"
	transit "github.com/aoeu/mta/transit_realtime"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

var re = regexp.MustCompile(`^\d+_L..N`)

type StopID string
type StopName string
type StopNames map[StopID]StopName

var URL string = "http://datamine.mta.info/mta_esi.php?key=%v&feed_id=2"
var key string = "TODO(aoeu): Implement a secure method for passing in an API key."
var stopsFileName string = "data/static/stops_l_train.json"
var stopNames StopNames

func GetNextLTrain() {
	b, err := ioutil.ReadFile(stopsFileName)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &stopNames)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf(URL, key)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := new(transit.FeedMessage)
	proto.Unmarshal(b, m)

	var entities []*transit.FeedEntity = m.GetEntity()

	for _, e := range entities {
		tu := e.GetTripUpdate()
		if tu == nil {
			// TODO(aoeu): What are these other messages that are not TripUpdates in the FeedMessage entities?
			continue
		}
		t := tu.GetTrip()
		stopTimeUpdates := tu.GetStopTimeUpdate()

		if !re.Match([]byte(t.GetTripId())) {
			continue
		}
		for _, stu := range stopTimeUpdates {
			id := stu.GetStopId()
			if id != "L13N" {
				continue
			}
			n := stopNames[StopID(id)]
			a := time.Unix(stu.GetArrival().GetTime(), 0)
			d := time.Unix(stu.GetDeparture().GetTime(), 0)
			next := time.Since(a)
			if next > 0 {
				continue
			}
			s := int(next.Seconds()) * -1
			m := int(s) / 60
			ss := int(s) - (m * 60)
			log.Printf("The next L train (%v) arrives at %v in %v minutes and %v seconds and departs at %v\n",
				t.GetTripId(), n, m, ss, d)
			return
		}
	}

}
