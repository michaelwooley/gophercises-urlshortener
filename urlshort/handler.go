package urlshort

import (
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

type pathsToUrlsYamlInterface []struct {
	// NOTE: Will not print out unless the elements are public!
	Path string
	URL  string
}

func (c *pathsToUrlsYamlInterface) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
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
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.Handler, error) {
	var pathsToUrlsYaml pathsToUrlsYamlInterface

	if err := pathsToUrlsYaml.Parse(yml); err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("Found %d url shortcuts", len(pathsToUrlsYaml)))

	pathsToUrls := make(map[string]string)
	for _, el := range pathsToUrlsYaml {
		pathsToUrls[el.Path] = el.URL
	}

	return MapHandler(pathsToUrls, "YAML", fallback), nil
}
