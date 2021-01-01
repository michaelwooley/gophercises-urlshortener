package urlshort

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, kind string, fallback http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := r.URL.Path
		val, ok := pathsToUrls[r.URL.Path]
		if !ok {
			log.Println(fmt.Sprintf("Prefix `%s` NOT found in %s. Redirecting to fallback.", prefix, kind))
			fallback.ServeHTTP(w, r)
		} else {
			log.Println(fmt.Sprintf("Prefix `%s` found in %s. Redirecting to: %s", prefix, kind, val))
			http.Redirect(w, r, val, http.StatusSeeOther)
		}
	})
}

type pathsToUrlsInterface []struct {
	// NOTE: Will not print out unless the elements are public!
	Path string
	URL  string
}

type urlMap map[string]string

func (c *pathsToUrlsInterface) info(kind string) {
	log.Printf("Found %d url shortcuts (%s)\n", len(*c), kind)

}

func (c *pathsToUrlsInterface) toMap() urlMap {
	out := make(urlMap)
	for _, cc := range *c {
		out[cc.Path] = cc.URL
	}
	return out
}

func (c *pathsToUrlsInterface) ParseYAML(data []byte) (urlMap, error) {
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	c.info("YAML")
	return c.toMap(), nil
}

func (c *pathsToUrlsInterface) ParseJSON(data []byte) (urlMap, error) {
	if err := json.Unmarshal(data, c); err != nil {
		return nil, err
	}
	c.info("JSON")
	return c.toMap(), nil
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
func YAMLHandler(data []byte, fallback http.Handler) (http.Handler, error) {
	var pathsToUrlsYaml pathsToUrlsInterface

	pathsToUrls, err := pathsToUrlsYaml.ParseYAML(data)
	if err != nil {
		return nil, err
	}

	return MapHandler(pathsToUrls, "YAML", fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
// 		[
//      	{
// 				"path": "/people",
//				"url": "people.com"
//			},
//      	{
// 				"path": "/example2",
//				"url": "example.com"
//			}
// 		]
//
// The only errors that can be returned all related to having
// invalid JSON data.
func JSONHandler(data []byte, fallback http.Handler) (http.Handler, error) {
	var pathsToUrlsJSON pathsToUrlsInterface

	pathsToUrls, err := pathsToUrlsJSON.ParseJSON(data)
	if err != nil {
		return nil, err
	}

	return MapHandler(pathsToUrls, "JSON", fallback), nil
}
