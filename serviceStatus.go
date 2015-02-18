package mta

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Service struct {
	Timestamp string `xml:"timestamp"`
	Subway    struct {
		Lines []Line `xml:"line"`
	} `xml:"subway"`
}

type Line struct {
	Name   string `xml:"name"`
	Status string `xml:"status"`
	Text   string `xml:"text"`
	Date   string `xml:"date"`
	Time   string `xml:"time"`
}

func GetLine(service Service, line string) (Line, error) {
	for _, l := range service.Subway.Lines {
		if strings.Contains(l.Name, line) {
			//l.Text = html.UnescapeString(l.Text)
			return l, nil
		}
	}
	return Line{}, errors.New("Line not found: " + line)
}

const serviceStatusURL = "http://mta.info/status/serviceStatus.txt"

// Gets the service status for all MTA subway lines.
func GetSubwayStatus(mtaurl ...string) (Service, error) {
	MTAURL := serviceStatusURL
	if len(mtaurl) > 0 {
		MTAURL = mtaurl[0]
	}
	body, err := getContents(MTAURL)
	if err != nil {
		return Service{}, err
	}
	sanitized := sanitizeInput(body)
	var service Service
	err = xml.Unmarshal(sanitized, &service)
	if err != nil {
		return Service{}, newError(err, "Error unmarshaling XML data:")
	}
	return service, err
}

func PrintSubwayStatus() {
	service, err := GetSubwayStatus()
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range service.Subway.Lines {
		fmt.Printf("%5s : %-20s\n", line.Name, line.Status)
	}
}

func getContents(source string) ([]byte, error) {
	MTAURL, err := url.Parse(source)
	if err != nil {
		return []byte{}, newError(err, "Error parsing URL:", MTAURL.String())
	}
	resp, err := http.Get(MTAURL.String())
	if err != nil {
		return []byte{}, newError(err, "Error getting URL:", MTAURL.String())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(&io.LimitedReader{resp.Body, 50000000})
	return body, err
}

func sanitizeInput(dirty []byte) []byte {
	clean := make([]byte, 0)
	for _, b := range dirty {
		if "U+0000" != fmt.Sprintf("%U", b) {
			clean = append(clean, b)
		}
	}
	return clean
}

func newError(err error, contexts ...string) error {
	var prefix string
	for _, context := range contexts {
		prefix += context + " "
	}
	return errors.New(prefix + err.Error())
}
