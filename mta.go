package mta

import (
	"encoding/json"
	"errors"
	"fmt"
	transit "github.com/aoeu/mta/transit_realtime"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

var re = regexp.MustCompile(`^\d+_L..N`)

type StopID string
type StopName string
type StopNames map[StopID]StopName

var URL string = "http://datamine.mta.info/mta_esi.php?key=%v&feed_id=2"
var key string = "TODO(aoeu): Implement a secure method for passing in an API key."
var stopsFileName string = "data/static/stops_l_train.json"

func GetFeedMessage(key string) (*transit.FeedMessage, error) {
	url := fmt.Sprintf(URL, key)
	resp, err := http.Get(url)
	if err != nil {
		e := errors.New(strings.Replace(err.Error(), key, "<key>", -1))
		return &transit.FeedMessage{}, e
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &transit.FeedMessage{}, err
	}
	m := &transit.FeedMessage{}
	err = proto.Unmarshal(b, m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func ReadStopNamesFile(stopsAsJSONFileName string) (StopNames, error) {
	s := make(StopNames, 0)
	b, err := ioutil.ReadFile(stopsAsJSONFileName)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(b, s)
	if err != nil {
		return s, err
	}
	return s, nil
}

type StopTime struct {
	TripID string
	StopName
	StopID
	Arrival   time.Time
	Departure time.Time
}

func NewStopTime(t *transit.TripDescriptor, s *transit.TripUpdate_StopTimeUpdate) *StopTime {
	stopID := StopID(s.GetStopId())
	return &StopTime{
		TripID:    t.GetTripId(),
		StopName:  stopNames[stopID],
		StopID:    stopID,
		Arrival:   time.Unix(s.GetArrival().GetTime(), 0),
		Departure: time.Unix(s.GetDeparture().GetTime(), 0),
	}
}

func (st StopTime) String() string {
	s := int(time.Since(st.Arrival).Seconds()) * -1
	m := int(s) / 60
	ss := int(s) - (m * 60)
	return fmt.Sprintf("The L train %v arrives at %v in %v minutes and %v seconds. Arrives at %v and departs at %v",
		st.TripID, st.StopName, m, ss, st.Arrival, st.Departure)

}

type StopTimes []StopTime

func (s *StopTimes) Len() int {
	return len(*s)
}

func (s *StopTimes) Less(i, j int) bool {
	si := (*s)[i]
	sj := (*s)[j]
	return si.Arrival.Before(sj.Arrival)
}

func (s *StopTimes) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func GetStopTimes(key string) ([]StopTime, error) {
	m, err := GetFeedMessage(key)
	if err != nil {
		return make([]StopTime, 0), err
	}
	var entities []*transit.FeedEntity = m.GetEntity() // TODO(aoeu): Can protoc name the method GetEntities()?
	stopTimes := make([]StopTime, 0)
	for _, e := range entities {
		var tripUpdate *transit.TripUpdate = e.GetTripUpdate()
		if tripUpdate == nil { // TODO(aoeu): Can FeedEntities have a type that we can switch on?
			continue
		}
		var stopTimeUpdates []*transit.TripUpdate_StopTimeUpdate = tripUpdate.GetStopTimeUpdate()
		var tripDescriptor *transit.TripDescriptor = tripUpdate.GetTrip()
		for _, stu := range stopTimeUpdates {
			st := NewStopTime(tripDescriptor, stu)
			stopTimes = append(stopTimes, *st)
		}
	}
	return stopTimes, nil
}

func GetNextMontroseLTrains(key string) (StopTimes, error) {
	st, err := GetStopTimes(key)
	if err != nil {
		return StopTimes{}, err
	}
	lst := make(StopTimes, 0)
	for _, s := range st {
		// TODO(aoeu): Why does it appear that bogus timestamps are on the data? Even other sites reflect this.
		if !re.Match([]byte(s.TripID)) { // TODO(aoeu): Assert that this check isn't needed.
			// log.Println("trip ID", s.TripID)
			continue
		}
		if s.StopID != "L13N" {
			// log.Println("stop ID", s.StopID)
			continue
		}
		if t := time.Since(s.Arrival); t > 0 {
			// log.Println("arrival time happened", t)
			continue
		}
		lst = append(lst, s)
	}
	sort.Sort(&lst)
	return lst, nil
}
