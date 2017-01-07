package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/aoeu/mta"
)

var port = flag.String("port", ":8080", "port on which to serve")
var lineNames = []string{
	"1", "2", "3",
	"4", "5", "6",
	"7",
	"A", "C", "E",
	"B", "D", "F", "M",
	"G",
	"J", "Z",
	"L",
	"N", "Q", "R",
	"S",
	"SIR"}

func main() {
	flag.Parse()
	service, err := mta.GetSubwayStatus()
	if err != nil {
		log.Fatalf("Problem getting subway status: %v", err)
	}
	allLineNames := make([]string, 0)
	for _, lineName := range lineNames {
		lineName = strings.ToUpper(lineName)
		handler := serveLine(service, lineName)
		http.HandleFunc("/"+lineName, handler)
		allLineNames = append(allLineNames, lineName)
		if l := strings.ToLower(lineName); l != lineName {
			http.HandleFunc("/"+l, handler)
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
var singleLineTempl = `
<html>
	<head>
		<title>{{.Name}} - MTA Service Status</title>
	</head>
	<body>
		<p>{{.Status}}</p>
		<p>{{.Text}}</p>
	</body>
</html>
`

func upper(s string) string {
	return strings.ToUpper(s)
}

func serveAllLines(lineNames []string) func(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{"upper": upper}
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
		templName := "train_line"
		templ, err := template.New(templName).Parse(singleLineTempl)
		if err != nil {
			log.Fatalf("Failed to parse template %v with error %v", templName, err)
		}
		line, err := mta.GetLine(service, lineName)
		if err != nil {
			http.Error(w, "Unknown line.", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Programmer error.", http.StatusInternalServerError)
			return
		}
		if err := templ.Execute(w, line); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
