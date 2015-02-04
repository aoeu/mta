package mta

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Service struct {
	Timestamp string `xml:"timestamp"`
	Subway    struct {
		Line []struct {
			Name   string `xml:"name"`
			Status string `xml:"status"`
			Text   string `xml:"text"`
			Date   string `xml:"date"`
			Time   string `xml:"time"`
		} `xml:"line"`
	} `xml:"subway"`
}

const serviceStatusURL = "http://mta.info/status/serviceStatus.txt"

func newError(fmtStr string, err error) error {
	return errors.New(fmt.Sprintf(fmtStr, err))
}

// Gets the service status for all MTA subway lines.
func GetSubwayStatus() (service Service, err error) {
	MTAURL, err := url.Parse(serviceStatusURL)
	if err != nil {
		return service, newError("Error parsing URL: %v", err)
	}
	resp, err := http.Get(MTAURL.String())
	if err != nil {
		return service, newError("Error getting URL: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return service, newError("Error reading response body: %v", err)
	}
	bodyBytes := make([]byte, 0)
	for _, rune := range body {
		if "U+0000" != fmt.Sprintf("%U", rune) {
			bodyBytes = append(bodyBytes, rune)
		}
	}
	err = xml.Unmarshal(bodyBytes, &service)
	if err != nil {
		return service, newError("Error unmarshaling XML data: %v", err)
	}
	return
}

func PrintSubwayStatus() {
	service, err := GetSubwayStatus()
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range service.Subway.Line {
		fmt.Printf("%5s : %-20s\n", line.Name, line.Status)
	}
}
