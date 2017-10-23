package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var version = "0.1.0"
var description = "Graphite Proxy instance - proxy graphite as a service"
var listen = "127.0.0.1"
var port = 9999
var graphiteAPIURL = "http://localhost:8888"

type MetricResponse struct {
	Text          string `json:"text"`
	Expandable    uint8  `json:"expandable"`
	Leaf          uint8  `json:"leaf"`
	ID            string `json:"id"`
	AllowChildren uint8  `json:"allowChildren"`
}

type MetricResponseError struct {
	Errors map[string]string `json:"errors"`
}

func metricErrors(errorMessage string) string {
	m := make(map[string]string, 1)
	m["query"] = errorMessage
	metricsErrors := MetricResponseError{
		Errors: m,
	}
	jsn, _ := json.Marshal(metricsErrors)
	log.Println("[ERROR]", errorMessage)
	return string(jsn)
}

func handlerRender(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{}
	login := req.Header.Get("login")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, metricErrors(err.Error()), 500)
		return
	}
	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := req.ParseForm(); err != nil {
		http.Error(w, metricErrors(err.Error()), 500)
		return
	}

	targets, ok := req.Form["target"]
	if !ok {
		http.Error(w, metricErrors("Targets not found"), 500)
		return
	}

	for _, target := range targets {
		// target is the metric + graphite functions
		// ex: aliasByMetric(toto.webserver.prod.toto-web.tail-apache2.counter-status_2xx)
		s := strings.Split(target, ".")
		s1 := strings.Split(s[0], "(")
		org := s1[len(s1)-1]

		// if user request render with target (metric) that is not
		// prefixed by his login, forbid request
		if !strings.HasPrefix(org, login) {
			log.Printf("401, Unauthorized (%s != %s)", org, login)
			http.Error(w, "401, Unauthorized", 401)
			return
		}
	}



	url := fmt.Sprintf("%s/render", graphiteAPIURL)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, metricErrors(err.Error()), 500)
		return
	}
	r.Header = req.Header
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, metricErrors(err.Error()), 500)
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, string(body))
	return

}

func handlerMetrics(w http.ResponseWriter, req *http.Request) {

	client := &http.Client{}
	login := req.Header.Get("login")

	// Get query value ex: /metrics/find?query=login.collectd.*
	// return "login.collectd.*"
	query := req.URL.Query().Get("query")

	// If query does not begin with login, return fake response
	// with only one choice: the login
	if !strings.HasPrefix(query, login) || query == "" {
		fakeResponses := make([]MetricResponse, 1)
		fakeResponses[0] = MetricResponse{
			Text:          login,
			Expandable:    1,
			Leaf:          0,
			ID:            login,
			AllowChildren: 1,
		}

		body, err := json.Marshal(fakeResponses)
		if err != nil {
			http.Error(w, metricErrors(err.Error()), 500)
			return
		}

		fmt.Fprintf(w, string(body))
		return
	}

	// Login has now been forced in the query
	url := fmt.Sprintf("%s/metrics/find?query=%s", graphiteAPIURL, query)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, metricErrors(err.Error()), 500)
		return
	}
	r.Header = req.Header

	// Do HTTP Request
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, metricErrors(err.Error()), 500)
		return
	}

	// Read Response and return back
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, metricErrors(err.Error()), 500)
		return
	}

	fmt.Fprintf(w, string(body))
	return
}

func handlerAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, _ := r.BasicAuth()
		if user == "" {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"Access Denied\"")
			http.Error(w, "401, Unauthorized", 401)
			return
		}
		r.Header.Set("login", user)
		fn(w, r)
	}
}

func handlerLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := "not authenticated"
		if r.Header.Get("login") != "" {
			login = r.Header.Get("login")
		}
		log.Printf("[INFO] %s %s %s%s (%s)", r.Method, r.Proto, r.Host, r.RequestURI, login)
		fn(w, r)
	}
}

func handlerVersion(w http.ResponseWriter, _ *http.Request) {
	var v = map[string]string{
		"version":     version,
		"description": description,
	}

	body, _ := json.Marshal(v)
	w.Write(body)
	return
}

func main() {
	listen := fmt.Sprintf("%s:%d", listen, port)

	http.HandleFunc("/metrics/find", handlerAuth(handlerLog(handlerMetrics)))
	http.HandleFunc("/render", handlerAuth(handlerLog(handlerRender)))
	http.HandleFunc("/version", handlerLog(handlerVersion))
	log.Fatal(http.ListenAndServe(listen, nil))

}
