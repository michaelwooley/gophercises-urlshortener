package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/michaelwooley/gophercises-urlshortener/urlshort"
)

// var yaml = `
// - path: /urlshort
//   url: https://github.com/gophercises/urlshort
// - path: /urlshort-final
//   url: https://github.com/gophercises/urlshort/tree/solution`

var jsonExample = `
[
	{
		"path": "/json",
		"url": "https://www.json.org/json-en.html"
	},
	{
		"path": "/geojson",
		"url": "https://tools.ietf.org/html/rfc7946"
	}
]
`

var yamlFilename string

func init() {
	const defaultYAML = "default-yaml.yml"

	flag.StringVar(&yamlFilename, "yaml-filename", defaultYAML, "Load this yaml file for use in ")
	flag.StringVar(&yamlFilename, "f", defaultYAML, "Shorthand for `-yaml-filename`")
}

func main() {
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, "default", mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yaml := readFromFile(yamlFilename)
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler([]byte(jsonExample), yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readFromFile(filename string) []byte {
	log.Printf("Reading file `%s`\n", filename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
