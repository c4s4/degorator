package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Process(writer http.ResponseWriter, request *http.Request, operation Operation) {
	log.Println(request.URL.String())
	parameters := request.URL.Query()
	// check request URL parameters
	for name, values := range parameters {
		parameter, ok := operation.Parameters[name]
		// if parameter is unknown, return 400 (Bad Request)
		if !ok {
			WriteResponse(writer, http.StatusBadRequest, "Unknown parameter %s", name)
			return
		}
		// check that parameter values match regexp, or return 400 (Bad Request)
		for _, value := range values {
			if !parameter.Compiled.MatchString(value) {
				WriteResponse(writer, http.StatusBadRequest, "Parameter %s doesn't match regexp %s", name, parameter.Regexp)
				return
			}
		}
	}
	// check for missing required parameters
	for name, parameter := range operation.Parameters {
		if !parameter.Optional {
			if _, ok := request.URL.Query()[name]; !ok {
				WriteResponse(writer, http.StatusBadRequest, "Missing mandatory parameter %s", name)
				return
			}
		}
	}
	// forward request to target
	url, err := url.Parse(config.Target.Host + operation.Target.Path)
	if err != nil {
		WriteResponse(writer, http.StatusInternalServerError, "Error building target URL: %v", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(writer, request)
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	for _, operation := range config.Operations {
		if request.Method == operation.Method && request.URL.Path == operation.Path {
			Process(writer, request, operation)
			return
		}
	}
	http.NotFound(writer, request)
}
