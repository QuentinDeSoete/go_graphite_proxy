package main

import (
	b64 "encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"bytes"
	"fmt"
)

var VERSION = "0.1.0"
var MESSAGE = "Graphite Proxy instance - proxy graphite as a service"

const (
	CONN_ROUTE = "/"
	CONN_PORT = ":9999"
)

func DecryptBasicAuth (sEnc string) (sDec []string) {
	s := strings.Split(sEnc, " ")
	temp, _ := b64.StdEncoding.DecodeString(s[1])
	sDec = strings.Split(string(temp), ":")
	return
}

func render(req *http.Request, login string, rw http.ResponseWriter, body []byte) {
	s := strings.Split(req.FormValue("target"), ".")
	s1 := strings.Split(s[0], "(")
	org := s1[len(s1)-1]
	if login == org {
		client := &http.Client{}
		r, err := http.NewRequest("POST", "http://localhost:8888/render", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		r.Header = req.Header
		resp, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		body1, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(rw, string(body1))
	}
}

func metric(req *http.Request, login string, rw http.ResponseWriter, body []byte) {
	s := strings.Split(req.RequestURI, "=")
	if s[1] == "*" {
		s1 := strings.Replace(req.RequestURI, "*", login, -1)
		client := &http.Client{}
		r, err := http.NewRequest("GET", fmt.Sprint("http://localhost:8888", s1), bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		r.Header = req.Header
		resp, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		body1, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(rw, string(body1))
	} else {
		s1 := strings.Split(s[1], ".")
		s2 := login
		for i := 1; i < len(s1); i++ {
			s2 = fmt.Sprint(s2, ".")
			s2 = fmt.Sprint(s2, s1[i])
		}
		client := &http.Client{}
		r, err := http.NewRequest("GET", fmt.Sprint("http://localhost:8888", fmt.Sprint(fmt.Sprint(s[0], "="), s2)), bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		r.Header = req.Header
		resp, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		body1, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(rw, string(body1))
	}
}

func test(rw http.ResponseWriter, req *http.Request) {
	if len(req.Header.Get("Authorization")) != 0 {
		body, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		login := DecryptBasicAuth(req.Header.Get("Authorization"))
		if req.URL.String() == "/render" {
			render(req, login[0], rw, body)
		} else {
			metric(req, login[0], rw, body)
		}
	} else {
		fmt.Fprintf(rw, "Unauthorized !")
	}
}

func main() {
	http.HandleFunc(CONN_ROUTE, test)
	http.ListenAndServe(CONN_PORT, nil)
}
