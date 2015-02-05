package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aoeu/mta"
)

var port = flag.String("port", ":8080", "port on which to serve")

func main() {
	flag.Parse()
	service, err := mta.GetSubwayStatus()
	if err != nil {
		log.Fatalf("Problem getting subway status: %v", err)
	}
	for _, lineName := range []string{"1", "2", "3", "4", "5", "6", "7", "A", "C", "E", "B", "D", "F", "M", "G", "J", "Z", "L", "N", "Q", "R", "S", "SIR"} {
		handler := serve(service, lineName)
		http.HandleFunc("/"+lineName, handler)
		if strings.ToUpper(lineName) != strings.ToLower(lineName) {
			http.HandleFunc("/"+strings.ToLower(lineName), handler)
		}
	}
	http.ListenAndServe(*port, nil)
}

func serve(service mta.Service, lineName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		line, err := mta.GetLine(service, lineName)
		if err != nil {
			http.Error(w, "Unknown line.", http.StatusNotFound)
			return
		}
		marshalled, err := json.Marshal(line)
		if err != nil {
			http.Error(w, "Programmer error.", http.StatusInternalServerError)
			return
		}
		//fmt.Fprint(w, "<!DOCTYPE html>"+string(marshalled))
		fmt.Fprint(w, begin+"line = "+string(marshalled)+"; \ndocument.write(line.Status); \nif (line.Text.length > 0) { \n\tdocument.write(line.Text); \n} else { \n\tdocument.Write(line.Status);\n}"+end)
	}
}

const begin = `
<script language="javascript">`
const end = `;
</script>
`
