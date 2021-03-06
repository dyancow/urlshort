package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//literally return a HandlerFunc
	return func(writer http.ResponseWriter, request *http.Request) {
		url := pathsToUrls[request.URL.Path]
		if len(url) != 0 {
			http.Redirect(writer, request, url, http.StatusFound)
			return
		}
		//otherwise use the fallback handler
		fallback.ServeHTTP(writer, request)
	}
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	//unmarshall the yml byte array into a map
	var yaml_arr []pathURL //an array of custom pathURL objects
	err := yaml.Unmarshal(yml, &yaml_arr)
	if err != nil {
		return nil, err
	}

	//now make yamls into maps
	yaml_map := make(map[string]string)
	for _, yaml_item := range yaml_arr { //kind like a python enumerate
		yaml_map[yaml_item.Path] = yaml_item.URL
	}

	//use MapHandler now that yaml inputs are in the same structure.
	return MapHandler(yaml_map, fallback), nil
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
