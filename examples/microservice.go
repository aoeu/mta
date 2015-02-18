package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
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
	allLineNames := make([]string, 0)
	for _, lineName := range []string{"1", "2", "3", "4", "5", "6", "7", "A", "C", "E", "B", "D", "F", "M", "G", "J", "Z", "L", "N", "Q", "R", "S", "SIR"} {
		handler := serveLine(service, lineName)
		http.HandleFunc("/"+lineName, handler)
		if strings.ToUpper(lineName) != strings.ToLower(lineName) {
			lineName := strings.ToLower(lineName)
			allLineNames = append(allLineNames, lineName)
			http.HandleFunc("/"+lineName, handler)
		}
	}
	http.HandleFunc("/", serveAllLines(allLineNames))
	http.ListenAndServe(*port, nil)
}

var allLinesTempl = `
<html>
	<head>
		<title>MTA Service Status</title>
	</head>
	<body>
	{{range $i, $lineName := .}}
		<a href="/{{.}}">{{upper .}}</a><br/>
	{{end}}
	</body>
</html>
`
func upper(s string) string {
	return strings.ToUpper(s)
}

func serveAllLines(lineNames []string) func(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{"upper" : upper}
	templName := "all_train_lines"
	templ, err := template.New(templName).Funcs(funcMap).Parse(allLinesTempl)
	if err != nil {
		log.Fatalf("Failed to parse template %v with error %v", templName, err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Execute(w, lineNames)
	}
}

func serveLine(service mta.Service, lineName string) func(w http.ResponseWriter, r *http.Request) {
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
